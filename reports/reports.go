package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/vlatan/postmark-reports/common"
)

type PostmarkReportsResponse struct {
	Entries []common.PostmarkReportInfo `json:"entries"`
	Meta    PostmarkMeta                `json:"meta"`
}

type PostmarkMeta struct {
	Next    any `json:"next"`
	NextURL any `json:"next_url"`
	Total   int `json:"total"`
}

func main() {
	godotenv.Load()
	saveReports("2024-11-18", "2024-11-25")
}

// https://dmarc.postmarkapp.com/api/#list-dmarc-reports
func saveReports(fromDate, toDate string) {
	client := http.Client{}
	url := fmt.Sprintf("%s/records/my/reports?from_date=%s&to_date=%s", common.DOMAIN, fromDate, toDate)
	reports := common.GetData(&client, url)
	var data PostmarkReportsResponse
	err := json.Unmarshal(reports, &data)
	common.Crash(err)

	result := []common.PostmarkReportInfo{}
	result = append(result, data.Entries...)

	for data.Meta.Next != nil {
		url = fmt.Sprintf("%s%s", common.DOMAIN, data.Meta.NextURL)
		reports = common.GetData(&client, url)
		err := json.Unmarshal(reports, &data)
		common.Crash(err)
		result = append(result, data.Entries...)
	}

	file, err := json.MarshalIndent(result, "", "\t")
	common.Crash(err)
	err = os.WriteFile("reports.json", file, 0644)
	common.Crash(err)
}
