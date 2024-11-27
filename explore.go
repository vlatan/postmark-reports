package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/vlatan/postmark-reports/common"
)

func main() {

	details, err := os.ReadFile("records.json")
	common.Crash(err)

	var data []common.PostmarkRecord
	err = json.Unmarshal(details, &data)
	common.Crash(err)

	total, failed, passed := 0, 0, 0
	amazonTotal, amazonPassed, amazonFailed := 0, 0, 0
	googleTotal, googlePassed, googleFailed := 0, 0, 0

	for _, record := range data {

		total += record.Count
		SPF := record.PolicyEvaluatedSpf == "pass"
		DKIM := record.PolicyEvaluatedDkim == "pass"
		amazon := strings.HasSuffix(record.HostName, "amazonses.com.")
		google := strings.HasSuffix(record.HostName, "google.com.")

		switch SPF || DKIM {
		case true:
			passed += record.Count
			if amazon {
				amazonTotal += record.Count
				amazonPassed += record.Count
			} else if google {
				googleTotal += record.Count
				googlePassed += record.Count
			}
		case false:
			failed += record.Count
			if amazon {
				amazonTotal += record.Count
				amazonFailed += record.Count
			} else if google {
				googleTotal += record.Count
				googleFailed += record.Count
			}
		}
	}

	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("Total Messages Sent:", total)
	fmt.Printf("SPF or DKIM aligned: %d (%.0f%%)\n", passed, Percent(passed, total))
	fmt.Printf("SPF and DKIM not aligned: %d (%.0f%%)\n", failed, Percent(failed, total))
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("Amazon Total Messages Sent:", amazonTotal)
	fmt.Printf("SPF or DKIM aligned: %d (%.0f%%)\n", amazonPassed, Percent(amazonPassed, amazonTotal))
	fmt.Printf("SPF and DKIM not aligned: %d (%.0f%%)\n", amazonFailed, Percent(amazonFailed, amazonTotal))
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("Google Total Messages Sent:", googleTotal)
	fmt.Printf("SPF or DKIM aligned: %d (%.0f%%)\n", googlePassed, Percent(googlePassed, googleTotal))
	fmt.Printf("SPF and DKIM not aligned:: %d (%.0f%%)\n", googleFailed, Percent(googleFailed, googleTotal))
	fmt.Println(strings.Repeat("-", 40))
}

func Percent(fraction, total int) float32 {
	return float32(fraction) / float32(total) * 100
}
