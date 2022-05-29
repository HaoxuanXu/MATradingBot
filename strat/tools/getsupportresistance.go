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
	var barClosePriceData []float64
	inputBarData := GetBarSubArray(barData[1:], time.Now(), 20)
	for _, bar := range inputBarData {
		barClosePriceData = append(barClosePriceData, bar.Close)
	}
	if metric == constants.SUPPORT {
		result, err = stats.Min(barClosePriceData)
		if err != nil {
			log.Panicf("support level calculation failed: %v", err)
		}
	} else if metric == constants.RESISTANCE {
		result, err = stats.Max(barClosePriceData)
		if err != nil {
			log.Panicf("resistance level calculation failed: %v", err)
		}
	}
	return result
}