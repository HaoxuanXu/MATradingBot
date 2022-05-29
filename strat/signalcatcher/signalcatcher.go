package signalcatcher

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.CloseData.CurrMAClose > model.CloseData.CurrMA20Close &&
		model.CloseData.CurrMAClose > model.CloseData.MAResistance &&
		time.Until(broker.Clock.NextClose) > time.Hour {
		return true
	}
	return false
}

func CanEnterShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.CloseData.CurrMAClose < model.CloseData.CurrMA20Close &&
		model.CloseData.CurrMAClose < model.CloseData.MASupport &&
		time.Until(broker.Clock.NextClose) > time.Hour {
		return true
	}
	return false
}
