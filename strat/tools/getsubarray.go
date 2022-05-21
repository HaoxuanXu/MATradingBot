package tools

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/util"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func GetBarSubArray(rawArray []marketdata.Bar, endTime time.Time, days int) []marketdata.Bar {
	var result []marketdata.Bar
	startTime := util.GetStartTime(endTime, days)

	for _, bar := range rawArray {
		if bar.Timestamp.Before(startTime) {
			break
		} else if bar.Timestamp.Before(endTime) {
			result = append(result, bar)
		}
	}

	return result
}
