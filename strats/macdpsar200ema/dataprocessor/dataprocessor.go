package dataprocessor

import (
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/cinar/indicator"
)

func ProcessBarData(model *model.DataModel, data *model.TotalBarData) bool {

	// update the bars
	model.Signal.Bars = data.StockBarData[model.Symbol]
	model.Signal.Quote = data.StockQuoteData[model.Symbol]
	if model.Signal.Bars[len(model.Signal.Bars)-1].Timestamp != model.CurrentBarTimestamp {
		// there is an update in the data, we then proceed the following processes
		// retrieve close data from bar slice
		var closeBars []float64
		var highBars []float64
		var lowBars []float64
		var volumeBars []int64
		for _, bar := range data.StockBarData[model.Symbol] {
			closeBars = append(closeBars, bar.Close)
			highBars = append(highBars, bar.High)
			lowBars = append(lowBars, bar.Low)
			volumeBars = append(volumeBars, int64(bar.Volume))
		}

		// calculate the current 200 period EMA value
		ema200Period := indicator.Ema(200, closeBars)
		model.Signal.EMA200Periods = ema200Period

		// calculate the current parabolic sar
		parSarVals, _ := indicator.ParabolicSar(highBars, lowBars, closeBars)
		model.Signal.ParabolicSars = parSarVals
		// calculate the current MACD values (MACD line, MACD signal line)
		macd, macdSignal := indicator.Macd(closeBars)
		model.Signal.Macds = macd
		model.Signal.MacdSignals = macdSignal

		// calculate Money Flow Index
		model.Signal.MoneyFlowIndex = indicator.DefaultMoneyFlowIndex(highBars, lowBars, closeBars, volumeBars)

		model.CurrentBarTimestamp = model.Signal.Bars[len(model.Signal.Bars)-1].Timestamp
		return true
	}

	return false

}
