package tools

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/strat/constants"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/montanaflynn/stats"
)

func CalcSupportResistance(barData []marketdata.Bar, metric constants.Metric) float64 {
	// we assume that the first value of the barData is the current value
	var result float64
	var err error
	var barHighPriceData []float64
	var barLowPriceData []float64
	inputBarData := GetBarSubArray(barData[1:], time.Now(), 20)
	for _, bar := range inputBarData {
		barHighPriceData = append(barHighPriceData, bar.High)
		barLowPriceData = append(barLowPriceData, bar.Low)
	}
	if metric == constants.SUPPORT {
		result, err = stats.Min(barLowPriceData)
		if err != nil {
			log.Fatalf("support level calculation failed: %v", err)
		}
	} else if metric == constants.RESISTANCE {
		result, err = stats.Max(barHighPriceData)
		if err != nil {
			log.Fatalf("resistance level calculation failed: %v", err)
		}
	}
	return result
}
