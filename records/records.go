package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/vlatan/postmark-reports/common"
)

type PostmarkRecords struct {
	CreatedAt           time.Time               `json:"created_at"`
	DateRangeBegin      time.Time               `json:"date_range_begin"`
	DateRangeEnd        time.Time               `json:"date_range_end"`
	DkimAlignmentMode   string                  `json:"dkim_alignment_mode"`
	Domain              string                  `json:"domain"`
	DomainPolicy        string                  `json:"domain_policy"`
	Email               string                  `json:"email"`
	ExternalID          string                  `json:"external_id"`
	ExtraContactInfo    string                  `json:"extra_contact_info"`
	FilteringPercentage int                     `json:"filtering_percentage"`
	ID                  int                     `json:"id"`
	OrganizationName    string                  `json:"organization_name"`
	RawSize             int                     `json:"raw_size"`
	RecordID            int                     `json:"record_id"`
	Records             []common.PostmarkRecord `json:"records"`
	SourceURI           string                  `json:"source_uri"`
	SpfAlignmentMode    string                  `json:"spf_alignment_mode"`
	SubdomainPolicy     string                  `json:"subdomain_policy"`
}

type Result struct {
	mutex sync.Mutex
	value []common.PostmarkRecord
}

type Job struct {
	client *http.Client
	report common.PostmarkReportInfo
}

func main() {
	godotenv.Load()
	reports, err := os.ReadFile("reports.json")
	common.Crash(err)
	compileReportsDetails(reports)
}

func (r *Result) Extend(records []common.PostmarkRecord) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.value = append(r.value, records...)
}

// https://dmarc.postmarkapp.com/api/#get-a-report-by-id
func getReportDetails(client *http.Client, id int) []byte {
	url := fmt.Sprintf("%s/records/my/reports/%d", common.DOMAIN, id)
	return common.GetData(client, url)
}

// Process jobs from the 'jobs' channel and update the result
func worker(jobs chan Job, result *Result, wg *sync.WaitGroup) {
	var recordsData PostmarkRecords
	for job := range jobs {
		records := getReportDetails(job.client, job.report.Id)
		err := json.Unmarshal(records, &recordsData)
		common.Crash(err)
		result.Extend(recordsData.Records)
		wg.Done()
	}
}

// Run workers, process the jobs (get records from reports)
// and combine the results into one file.
func compileReportsDetails(reports []byte) {
	numWorkers, err := strconv.Atoi(os.Getenv("WORKERS"))
	common.Crash(err)

	var data []common.PostmarkReportInfo
	err = json.Unmarshal(reports, &data)
	common.Crash(err)

	var result Result
	var wg sync.WaitGroup
	jobs := make(chan Job, len(data))

	// Start the workers in a separate goroutines.
	// Each will pop a job from the 'jobs' channel
	// and process it concurrently.
	for w := 0; w < numWorkers; w++ {
		go worker(jobs, &result, &wg)
	}

	// Queue up the jobs
	client := http.Client{}
	for _, entry := range data {
		wg.Add(1)
		jobs <- Job{&client, entry}
	}

	// Wait for jobs to finish
	wg.Wait()

	// We don't even need to close the channel
	// because when wg.Wait() is done it means
	// all the records are already in the Result
	// and we can proceed with the file creation.
	close(jobs)

	// Create one file
	file, err := json.MarshalIndent(result.value, "", "\t")
	common.Crash(err)
	err = os.WriteFile("records.json", file, 0644)
	common.Crash(err)
}
