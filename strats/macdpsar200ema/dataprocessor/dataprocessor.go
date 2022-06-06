package dataprocessor

import (
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/cinar/indicator"
)

func ProcessBarData(model *model.DataModel, data *model.TotalBarData) bool {

	// update the current bar
	model.Signal.CurrentBar = data.BarData[model.Symbol][len(data.BarData[model.Symbol])-1]

	if model.Signal.CurrentBar.Timestamp != model.CurrentBarTimestamp {
		// there is an update in the data, we then proceed the following processes
		// retrieve close data from bar slice
		var closeBars []float64
		var highBars []float64
		var lowBars []float64
		for _, bar := range data.BarData[model.Symbol] {
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

		// calculate the current MACD values (MACD line, MACD signal line)
		macd, macdSignal := indicator.Macd(closeBars)
		model.Signal.CurrentMacd = macd[len(macd)-1]
		model.Signal.CurrentMacdSignal = macdSignal[len(macdSignal)-1]
		return true
	}
	return false

}
