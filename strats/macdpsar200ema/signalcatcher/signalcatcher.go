package signalcatcher

import (
	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Signal.CurrentBar.Close > model.Signal.CurrentEMA200Period && // the current price is above the 200 period EMA value
		model.Signal.CurrentParabolicSar < model.Signal.CurrentBar.Low &&
		model.Signal.PreviousParabolicSar > model.Signal.PreviousBar.High &&
		model.Signal.CurrentMacd > model.Signal.CurrentMacdSignal {
		return true
	}

	return false
}

func CanEnterShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Signal.CurrentBar.Close < model.Signal.CurrentEMA200Period && // the current price is below the 200 period EMA value
		model.Signal.CurrentParabolicSar > model.Signal.CurrentBar.High &&
		model.Signal.PreviousParabolicSar < model.Signal.PreviousBar.Low &&
		model.Signal.CurrentMacd < model.Signal.CurrentMacdSignal {
		return true
	}
	return false
}
