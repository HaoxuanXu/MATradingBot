package signalcatcher

import (
	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func pSarPrevDirection(pSars []float64, pBars []marketdata.Bar, currDirection string) bool {
	var resultFloats []float64
	if currDirection == "above" {
		var pBarFloats []float64

		for _, bar := range pBars {
			pBarFloats = append(pBarFloats, bar.Low)
		}

		for i := 0; i < len(pBarFloats); i++ {
			resultFloats = append(resultFloats, pBarFloats[i]-pSars[i])
		}

		for _, val := range resultFloats {
			if val > 0 {
				return true
			}
		}
	} else if currDirection == "below" {
		var pBarFloats []float64

		for _, bar := range pBars {
			pBarFloats = append(pBarFloats, bar.High)
		}

		for i := 0; i < len(pBarFloats); i++ {
			resultFloats = append(resultFloats, pBarFloats[i]-pSars[i])
		}

		for _, val := range resultFloats {
			if val < 0 {
				return true
			}
		}

	}

	return false

}

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Signal.Bars[len(model.Signal.Bars)-1].Close > model.Signal.EMA200Periods[len(model.Signal.EMA200Periods)-1] && // the current price is above the 200 period EMA value
		model.Signal.ParabolicSars[len(model.Signal.ParabolicSars)-1] < model.Signal.Quote.AskPrice &&
		pSarPrevDirection(model.Signal.ParabolicSars[len(model.Signal.ParabolicSars)-4:], model.Signal.Bars[len(model.Signal.Bars)-4:], "below") &&
		model.Signal.Macds[len(model.Signal.Macds)-1] > model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] {
		return true
	}

	return false
}

func CanEnterShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Signal.Bars[len(model.Signal.Bars)-1].Close < model.Signal.EMA200Periods[len(model.Signal.EMA200Periods)-1] && // the current price is below the 200 period EMA value
		model.Signal.ParabolicSars[len(model.Signal.ParabolicSars)-1] > model.Signal.Quote.BidPrice &&
		pSarPrevDirection(model.Signal.ParabolicSars[len(model.Signal.ParabolicSars)-4:], model.Signal.Bars[len(model.Signal.Bars)-4:], "above") &&
		model.Signal.Macds[len(model.Signal.Macds)-1] < model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] {
		return true
	}
	return false
}
