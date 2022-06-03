package tools

import (
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/montanaflynn/stats"
)

func CalcMovingAverage(barData []marketdata.Bar, endTime time.Time, periods int) float64 {
	var closeArray []float64

	barSubArray := GetBarSubArray(barData, endTime, periods)
	for _, bar := range barSubArray {
		closeArray = append(closeArray, bar.Close)
	}

	result, _ := stats.Mean(closeArray)
	result, _ = stats.Round(result, 3)
	return result
}
