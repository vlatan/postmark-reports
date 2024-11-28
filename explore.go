package main

import (
	"encoding/json"
	"os"

	"github.com/vlatan/postmark-reports/common"
)

type Count struct {
	Total          int `json:"messages_sent"`
	SPForDKIMPass  int `json:"spf_or_dkim_pass"`
	SPFandDKIMFail int `json:"spf_and_dkim_fail"`
	SPFPass        int `json:"spf_pass"`
	DKIMPass       int `json:"dkim_pass"`
}
type Domains map[string]Count

type Stats struct {
	Total   Count   `json:"total"`
	Domains Domains `json:"domains"`
}

func main() {

	details, err := os.ReadFile("records.json")
	common.Crash(err)

	var data []common.PostmarkRecord
	err = json.Unmarshal(details, &data)
	common.Crash(err)

	total, totalSPFPass, totalDKIMPass, totalSPForDKIMPass := 0, 0, 0, 0
	totalSPFandDKIMFail := 0
	domains := make(Domains)

	for _, record := range data {
		total += record.Count
		domain := record.TopPrivateDomainName
		if domain == "" || domain == "." {
			domain = "unresolved"
		}

		domainCount := domains[domain]
		domainCount.Total += record.Count

		countPass := 0
		if record.PolicyEvaluatedSpf == "pass" {
			totalSPFPass += record.Count
			domainCount.SPFPass += record.Count
			countPass++
		}

		if record.PolicyEvaluatedDkim == "pass" {
			totalDKIMPass += record.Count
			domainCount.DKIMPass += record.Count
			countPass++
		}

		switch countPass {
		case 0:
			totalSPFandDKIMFail += record.Count
			domainCount.SPFandDKIMFail += record.Count
		case 1, 2:
			totalSPForDKIMPass += record.Count
			domainCount.SPForDKIMPass += record.Count

		}

		domains[domain] = domainCount
	}

	stats := Stats{
		Total: Count{
			Total:          total,
			SPForDKIMPass:  totalSPForDKIMPass,
			SPFandDKIMFail: totalSPFandDKIMFail,
			SPFPass:        totalSPFPass,
			DKIMPass:       totalDKIMPass,
		},
		Domains: domains,
	}

	file, err := json.MarshalIndent(stats, "", "\t")
	common.Crash(err)
	err = os.WriteFile("stats.json", file, 0644)
	common.Crash(err)
}
