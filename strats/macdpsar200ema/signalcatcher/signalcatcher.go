package signalcatcher

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/HaoxuanXu/MATradingBot/util"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLong(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Signal.CurrentBar.Low > model.Signal.CurrentEMA200Period && // the current price is above the 200 period EMA value
		model.Signal.RecentParabolicSarDiff[len(model.Signal.RecentParabolicSarDiff)-1] < 0 &&
		util.ValueInArraySmaller(model.Signal.RecentParabolicSarDiff, 0) &&
		time.Until(broker.Clock.NextClose) < time.Hour {
		return true
	}
	return false
}

func CanEnterShort(model *model.DataModel, broker *api.AlpacaBroker) bool {
	if !model.Position.HasLongPosition && !model.Position.HasShortPosition &&
		model.Signal.CurrentBar.High < model.Signal.CurrentEMA200Period && // the current price is below the 200 period EMA value
		model.Signal.RecentParabolicSarDiff[len(model.Signal.RecentParabolicSarDiff)-1] > 0 &&
		util.ValueInArrayLarger(model.Signal.RecentParabolicSarDiff, 0) &&
		time.Until(broker.Clock.NextClose) < time.Hour {
		return true
	}
	return false
}
