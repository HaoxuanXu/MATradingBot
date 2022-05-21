package signalcatcher

import "github.com/HaoxuanXu/MATradingBot/strat/model"

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(conditions model.MAConditions) bool {
	if !conditions.HasLongPosition && !conditions.HasShortPosition &&
		conditions.IsMA20DaysRising && conditions.IsMA30DaysRising &&
		conditions.IsMA20AboveMA30 {
		return true
	}
	return false
}

func CanEnterShort(conditions model.MAConditions) bool {
	if !conditions.HasLongPosition && !conditions.HasShortPosition &&
		conditions.IsMA20DaysDropping && conditions.IsMA30DaysDropping &&
		conditions.IsMA20BelowMA30 {
		return true
	}
	return false
}

func CanCloseLong(conditions model.MAConditions) bool {
	if conditions.HasLongPosition && !conditions.HasShortPosition &&
		conditions.IsMA20DaysDropping {
		return true
	}
	return false
}

func CanCloseShort(conditions model.MAConditions) bool {
	if conditions.HasShortPosition && !conditions.HasLongPosition &&
		conditions.IsMA20DaysRising {
		return true
	}
	return false
}