package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/vlatan/postmark-reports/common"
)

type ReportRecords struct {
	CreatedAt           time.Time       `json:"created_at"`
	DateRangeBegin      time.Time       `json:"date_range_begin"`
	DateRangeEnd        time.Time       `json:"date_range_end"`
	DkimAlignmentMode   string          `json:"dkim_alignment_mode"`
	Domain              string          `json:"domain"`
	DomainPolicy        string          `json:"domain_policy"`
	Email               string          `json:"email"`
	ExternalID          string          `json:"external_id"`
	ExtraContactInfo    string          `json:"extra_contact_info"`
	FilteringPercentage int             `json:"filtering_percentage"`
	ID                  int             `json:"id"`
	OrganizationName    string          `json:"organization_name"`
	RawSize             int             `json:"raw_size"`
	RecordID            int             `json:"record_id"`
	Records             []common.Record `json:"records"`
	SourceURI           string          `json:"source_uri"`
	SpfAlignmentMode    string          `json:"spf_alignment_mode"`
	SubdomainPolicy     string          `json:"subdomain_policy"`
}

type Result struct {
	mutex sync.Mutex
	value []common.Record
}

func (r *Result) Extend(records []common.Record) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.value = append(r.value, records...)
}

func main() {
	godotenv.Load()
	reports, err := os.ReadFile("reports.json")
	common.Crash(err)
	compileReportsDetails(reports)
}

// https://dmarc.postmarkapp.com/api/#get-a-report-by-id
func getReportDetails(client *http.Client, id int) []byte {
	url := fmt.Sprintf("%s/records/my/reports/%d", common.DOMAIN, id)
	return common.GetData(client, url)
}

func compileReportsDetails(reports []byte) {
	var data []common.ReportEntry
	err := json.Unmarshal(reports, &data)
	common.Crash(err)

	client := http.Client{}
	var result Result
	var wg sync.WaitGroup

	for _, entry := range data {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var recordsData ReportRecords
			records := getReportDetails(&client, entry.Id)
			err = json.Unmarshal(records, &recordsData)
			common.Crash(err)
			result.Extend(recordsData.Records)
		}()
	}

	wg.Wait()
	file, err := json.MarshalIndent(result.value, "", "\t")
	common.Crash(err)
	err = os.WriteFile("records.json", file, 0644)
	common.Crash(err)
}
