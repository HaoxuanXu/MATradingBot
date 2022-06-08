package dataprocessor

import (
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/cinar/indicator"
)

func ProcessBarData(model *model.DataModel, data *model.TotalBarData, isCrypto bool) bool {

	// update the current bar
	if !isCrypto {
		model.Signal.CurrentBar = data.StockBarData[model.Symbol][len(data.StockBarData[model.Symbol])-1]
		model.Signal.PreviousBar = data.StockBarData[model.Symbol][len(data.StockBarData[model.Symbol])-2]
	} else {
		model.Signal.CurrentCryptoBar = data.CryptoBarData[model.Symbol][len(data.CryptoBarData[model.Symbol])-1]
		model.Signal.PreviousCryptoBar = data.CryptoBarData[model.Symbol][len(data.CryptoBarData[model.Symbol])-2]

	}

	if !isCrypto {
		if model.Signal.CurrentBar.Timestamp != model.CurrentBarTimestamp {
			// there is an update in the data, we then proceed the following processes
			// retrieve close data from bar slice
			var closeBars []float64
			var highBars []float64
			var lowBars []float64
			for _, bar := range data.StockBarData[model.Symbol] {
				closeBars = append(closeBars, bar.Close)
				highBars = append(highBars, bar.High)
				lowBars = append(lowBars, bar.Low)
			}

			// calculate the current 200 period EMA value
			ema200Period := indicator.Ema(200, closeBars)
			model.Signal.CurrentEMA200Period = ema200Period[len(ema200Period)-1]

			// calculate the current parabolic sar
			parSarVals, _ := indicator.ParabolicSar(highBars, lowBars, closeBars)
			model.Signal.CurrentParabolicSar = parSarVals[len(parSarVals)-1]
			model.Signal.PreviousParabolicSar = parSarVals[len(parSarVals)-2]
			// calculate the current MACD values (MACD line, MACD signal line)
			macd, macdSignal := indicator.Macd(closeBars)
			model.Signal.CurrentMacd = macd[len(macd)-1]
			model.Signal.CurrentMacdSignal = macdSignal[len(macdSignal)-1]

			model.CurrentBarTimestamp = model.Signal.CurrentBar.Timestamp
			return true
		}
	} else {
		// there is an update in the data, we then proceed the following processes
		// retrieve close data from bar slice
		var closeBars []float64
		var highBars []float64
		var lowBars []float64
		for _, bar := range data.CryptoBarData[model.Symbol] {
			closeBars = append(closeBars, bar.Close)
			highBars = append(highBars, bar.High)
			lowBars = append(lowBars, bar.Low)
		}

		// calculate the current 200 period EMA value
		ema200Period := indicator.Ema(200, closeBars)
		model.Signal.CurrentEMA200Period = ema200Period[len(ema200Period)-1]

		// calculate the current parabolic sar
		parSarVals, _ := indicator.ParabolicSar(highBars, lowBars, closeBars)
		model.Signal.CurrentParabolicSar = parSarVals[len(parSarVals)-1]
		model.Signal.PreviousParabolicSar = parSarVals[len(parSarVals)-2]
		// calculate the current MACD values (MACD line, MACD signal line)
		macd, macdSignal := indicator.Macd(closeBars)
		model.Signal.CurrentMacd = macd[len(macd)-1]
		model.Signal.CurrentMacdSignal = macdSignal[len(macdSignal)-1]

		model.CurrentBarTimestamp = model.Signal.CurrentCryptoBar.Timestamp
		return true
	}

	return false

}
