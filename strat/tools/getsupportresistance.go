package tools

import (
	"log"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/montanaflynn/stats"
)

func CalcSupport(barData []marketdata.Bar) float64 {
	// we assume that the first value of the barData is the current value
	var barClosePriceData []float64
	inputBarData := GetBarSubArray(barData[1:], time.Now(), 30)
	for _, bar := range inputBarData {
		barClosePriceData = append(barClosePriceData, bar.Close)
	}
	support, err := stats.Min(barClosePriceData)
	if err != nil {
		log.Panicf("support level calculation failed: %v", err)
	}
	return support
}

func CalcResistance(barData []marketdata.Bar) float64 {
	// we assume that the first value of the barData is the current value
	var barClosePriceData []float64
	inputBarData := barData[1:]
	for _, bar := range inputBarData {
		barClosePriceData = append(barClosePriceData, bar.Close)
	}
	resistance, err := stats.Max(barClosePriceData)
	if err != nil {
		log.Panicf("resistance level calculation failed: %v", err)
	}
	return resistance
}
