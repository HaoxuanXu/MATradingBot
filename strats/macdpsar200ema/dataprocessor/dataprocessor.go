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

		// calculate BXTrender
		Ema5DaysClose := indicator.Ema(5, closeBars)
		Ema20DaysClose := indicator.Ema(20, closeBars)
		Ema20DaysClose = Ema20DaysClose[len(Ema5DaysClose)-len(Ema20DaysClose):]
		var priceDiff []float64

		for i := 0; i < len(Ema5DaysClose); i++ {
			priceDiff = append(priceDiff, Ema5DaysClose[i]-Ema20DaysClose[i])
		}

		_, bxTrenderShortTerm := indicator.Rsi(Ema20DaysClose)
		model.Signal.BXTrenderShortTerm = bxTrenderShortTerm

		_, bxtrenderLongTerm := indicator.Rsi(priceDiff)
		model.Signal.BXTrenderLongTerm = bxtrenderLongTerm

		// calculate swing low swing high
		swingLow := indicator.Min(13, closeBars)
		swingHigh := indicator.Max(13, closeBars)
		model.Signal.SwingLow = swingLow
		model.Signal.SwingHigh = swingHigh
		var fibonacciLow []float64
		var fibonacciHigh []float64
		for i := 0; i < len(swingLow); i++ {
			fibonacciLow = append(fibonacciLow, swingLow[i]+(swingHigh[i]-swingLow[i])*0.382)
			fibonacciHigh = append(fibonacciHigh, swingHigh[i]-(swingHigh[i]-swingLow[i])*0.382)
		}
		model.Signal.FibonacciLow = fibonacciLow
		model.Signal.FibonacciHigh = fibonacciHigh

		// calculate stochastic
		k, d := indicator.StochasticOscillator(highBars, lowBars, closeBars)
		model.Signal.StochK = k
		model.Signal.StochD = d
		// determine overbought oversold
		if model.Signal.StochK[len(model.Signal.StochK)-1] < 20 && model.Signal.StochD[len(model.Signal.StochD)-1] < 20 {
			model.Signal.StochOversold = true
			model.Signal.StochOverbought = false
		} else if model.Signal.StochK[len(model.Signal.StochK)-1] > 80 && model.Signal.StochD[len(model.Signal.StochD)-1] > 80 {
			model.Signal.StochOversold = false
			model.Signal.StochOverbought = true
		}

		//

		// calculate RSI
		_, rsi := indicator.Rsi(closeBars)
		model.Signal.RSI = rsi

		// calculate trailing stop loss
		exitLong, exitShort := indicator.ChandelierExit(highBars, lowBars, closeBars)
		model.Signal.TrailingStopLossLong = exitLong[len(exitLong)-1]
		model.Signal.TrailingStopLossShort = exitShort[len(exitShort)-1]

		model.CurrentBarTimestamp = model.Signal.Bars[len(model.Signal.Bars)-1].Timestamp
		return true
	}

	return false

}
