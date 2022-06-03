package signalcatcher

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.CloseData.CurrMATrade > model.CloseData.CurrMA20Close &&
		// model.CloseData.CurrMATrade > model.CloseData.MAResistance &&
		model.CloseData.CurrMA20Close > model.CloseData.PrevMA20Close &&
		model.Trails.AppliedLongTrail > model.CloseData.CurrMATrade*0.001 &&
		model.CloseData.CurrMATrade-model.CloseData.PrevMATrade > 0 &&
		time.Until(broker.Clock.NextClose) > time.Hour {
		return true
	}
	return false
}

func CanEnterShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.CloseData.CurrMATrade < model.CloseData.CurrMA20Close &&
		// model.CloseData.CurrMATrade < model.CloseData.MASupport &&
		model.CloseData.CurrMA20Close < model.CloseData.PrevMA20Close &&
		model.Trails.AppliedShortTrail > model.CloseData.CurrMATrade*0.001 &&
		model.CloseData.CurrMATrade-model.CloseData.PrevMATrade < 0 &&
		time.Until(broker.Clock.NextClose) > time.Hour {
		return true
	}
	return false
}

func CanExitLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.CloseData.CurrMATrade > model.Position.FilledPrice &&
		time.Until(broker.Clock.NextClose) < time.Hour {
		return true
	}
	return false
}

func CanExitShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if model.Position.HasShortPosition && !model.Position.HasLongPosition &&
		model.CloseData.CurrMATrade < model.Position.FilledPrice &&
		time.Until(broker.Clock.NextClose) < time.Hour {
		return true
	}
	return false
}
