package signalcatcher

import (
	"github.com/HaoxuanXu/MATradingBot/strat/model"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Condition.IsMA20PeriodsRising &&
		model.Condition.IsMAAboveMA20 && model.CloseData.CurrMAClose > model.CloseData.MAResistance {
		return true
	}
	return false
}

func CanEnterShort(model *model.DataModel) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Condition.IsMA20PeriodsDropping &&
		model.Condition.IsMABelowMA20 && model.CloseData.CurrMAClose < model.CloseData.MASupport {
		return true
	}
	return false
}

func CanCloseLong(conditions model.MAConditions, position model.PositionData) bool {
	if position.HasLongPosition && !position.HasShortPosition &&
		conditions.IsMA20PeriodsDropping {
		return true
	}
	return false
}

func CanCloseShort(conditions model.MAConditions, position model.PositionData) bool {
	if position.HasShortPosition && !position.HasLongPosition &&
		conditions.IsMA20PeriodsRising {
		return true
	}
	return false
}
