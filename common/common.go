package common

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ReportEntry struct {
	Id               int       `json:"id"`
	Domain           string    `json:"domain"`
	ExternalId       string    `json:"external_id"`
	OrganizationName string    `json:"organization_name"`
	DateRangeBegin   time.Time `json:"date_range_begin"`
	DateRangeEnd     time.Time `json:"date_range_end"`
	CreatedAt        time.Time `json:"created_at"`
}

type Record struct {
	Count                        int       `json:"count"`
	CreatedAt                    time.Time `json:"created_at"`
	DkimDomain                   any       `json:"dkim_domain"`
	DkimResult                   any       `json:"dkim_result"`
	HeaderFrom                   string    `json:"header_from"`
	HostName                     string    `json:"host_name"`
	PolicyEvaluatedDisposition   string    `json:"policy_evaluated_disposition"`
	PolicyEvaluatedDkim          string    `json:"policy_evaluated_dkim"`
	PolicyEvaluatedReasonComment any       `json:"policy_evaluated_reason_comment"`
	PolicyEvaluatedReasonType    any       `json:"policy_evaluated_reason_type"`
	PolicyEvaluatedSpf           string    `json:"policy_evaluated_spf"`
	RecordID                     int       `json:"record_id"`
	ReportID                     int       `json:"report_id"`
	RowNum                       int       `json:"row_num"`
	SourceIP                     string    `json:"source_ip"`
	SourceIPVersion              int       `json:"source_ip_version"`
	SpfDomain                    string    `json:"spf_domain"`
	SpfResult                    string    `json:"spf_result"`
	TopPrivateDomainName         string    `json:"top_private_domain_name"`
}

const DOMAIN = "https://dmarc.postmarkapp.com"

func Crash(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetData(client *http.Client, url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	Crash(err)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Api-Token", os.Getenv("TOKEN"))

	resp, err := client.Do(req)
	Crash(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	Crash(err)

	return body
}
