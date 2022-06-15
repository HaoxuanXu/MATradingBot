package signalcatcher

import (
	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Signal.Bars[len(model.Signal.Bars)-1].Close > model.Signal.EMA200Periods[len(model.Signal.EMA200Periods)-1] && // the current price is above the 200 period EMA value
		model.Signal.Macds[len(model.Signal.Macds)-1] > model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] &&
		model.Signal.Macds[len(model.Signal.Macds)-1] < 0 &&
		model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] < 0 &&
		model.Signal.StochOversold &&
		!model.Signal.StockOverbought &&
		model.Signal.RSI[len(model.Signal.RSI)-1] > 50 {
		return true
	}

	return false
}

func CanEnterShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Signal.Bars[len(model.Signal.Bars)-1].Close < model.Signal.EMA200Periods[len(model.Signal.EMA200Periods)-1] && // the current price is below the 200 period EMA value
		model.Signal.Macds[len(model.Signal.Macds)-1] < model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] &&
		model.Signal.Macds[len(model.Signal.Macds)-1] > 0 &&
		model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] > 0 &&
		!model.Signal.StochOversold &&
		model.Signal.StockOverbought &&
		model.Signal.RSI[len(model.Signal.RSI)-1] < 50 {
		return true
	}
	return false
}
