package signalcatcher

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.CloseData.CurrMAAsk > model.CloseData.CurrMA20Close &&
		model.CloseData.CurrMAAsk > model.CloseData.MAResistance &&
		model.CloseData.CurrMA20Close > model.CloseData.PrevMA20Close &&
		model.Trails.AppliedLongTrail > model.CloseData.CurrMAAsk*0.001 &&
		time.Until(broker.Clock.NextClose) > time.Hour {
		return true
	}
	return false
}

func CanEnterShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.CloseData.CurrMABid < model.CloseData.CurrMA20Close &&
		model.CloseData.CurrMABid < model.CloseData.MASupport &&
		model.CloseData.CurrMA20Close < model.CloseData.PrevMA20Close &&
		model.Trails.AppliedShortTrail > model.CloseData.CurrMABid*0.001 &&
		time.Until(broker.Clock.NextClose) > time.Hour {
		return true
	}
	return false
}

func CanExitLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.CloseData.CurrMABid > model.Position.FilledPrice &&
		time.Until(broker.Clock.NextClose) < time.Hour {
		return true
	}
	return false
}

func CanExitShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if model.Position.HasShortPosition && !model.Position.HasLongPosition &&
		model.CloseData.CurrMAAsk < model.Position.FilledPrice &&
		time.Until(broker.Clock.NextClose) < time.Hour {
		return true
	}
	return false
}
