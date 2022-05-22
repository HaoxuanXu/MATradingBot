package signalcatcher

import (
	"github.com/HaoxuanXu/MATradingBot/strat/model"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Condition.IsMA20DaysRising && model.Condition.IsMA30DaysRising &&
		model.Condition.IsMA20AboveMA30 {
		return true
	}
	return false
}

func CanEnterShort(model *model.DataModel) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Condition.IsMA20DaysDropping && model.Condition.IsMA30DaysDropping &&
		model.Condition.IsMA20BelowMA30 {
		return true
	}
	return false
}

func CanCloseLong(conditions model.MAConditions, position model.PositionData) bool {
	if position.HasLongPosition && !position.HasShortPosition &&
		conditions.IsMA20DaysDropping {
		return true
	}
	return false
}

func CanCloseShort(conditions model.MAConditions, position model.PositionData) bool {
	if position.HasShortPosition && !position.HasLongPosition &&
		conditions.IsMA20DaysRising {
		return true
	}
	return false
}
