package main

import (
	"encoding/json"
	"os"

	"github.com/vlatan/postmark-reports/common"
)

type Count struct {
	Total    int `json:"total_messages_sent"`
	SPFPass  int `json:"spf_pass"`
	DKIMPass int `json:"dkim_pass"`
}
type Domains map[string]Count

type Stats struct {
	Total    int     `json:"total_messages_sent"`
	SPFPass  int     `json:"spf_pass"`
	DKIMPass int     `json:"dkim_pass"`
	Domains  Domains `json:"domains"`
}

func main() {

	details, err := os.ReadFile("records.json")
	common.Crash(err)

	var data []common.PostmarkRecord
	err = json.Unmarshal(details, &data)
	common.Crash(err)

	total, totalSPFPass, totalDKIMPass := 0, 0, 0
	domains := make(Domains)

	for _, record := range data {
		total += record.Count
		domain := record.TopPrivateDomainName
		if domain == "" || domain == "." {
			domain = "unresolved"
		}

		domainCount := domains[domain]
		domainCount.Total += record.Count

		if record.PolicyEvaluatedSpf == "pass" {
			totalSPFPass += record.Count
			domainCount.SPFPass += record.Count
		}

		if record.PolicyEvaluatedDkim == "pass" {
			totalDKIMPass += record.Count
			domainCount.DKIMPass += record.Count
		}

		domains[domain] = domainCount
	}

	stats := Stats{
		Total:    total,
		SPFPass:  totalSPFPass,
		DKIMPass: totalDKIMPass,
		Domains:  domains,
	}

	file, err := json.MarshalIndent(stats, "", "\t")
	common.Crash(err)
	err = os.WriteFile("stats.json", file, 0644)
	common.Crash(err)
}
