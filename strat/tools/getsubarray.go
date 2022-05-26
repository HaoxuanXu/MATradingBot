package tools

import (
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func GetBarSubArray(rawArray []marketdata.Bar, endTime time.Time, periods int) []marketdata.Bar {
	var result []marketdata.Bar

	for _, bar := range rawArray {
		if len(result) < periods {
			if bar.Timestamp.Before(endTime) {
				result = append(result, bar)
			}
		}
	}

	return result
}
