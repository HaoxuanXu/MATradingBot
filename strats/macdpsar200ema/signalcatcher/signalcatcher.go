package signalcatcher

import (
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
)

// In order to go long, 20MA has to be above  30MA and both MAs have to be rising
func CanEnterLongFastPeriod(model *model.DataModel) bool {
	if model.Signal.Macds[len(model.Signal.Macds)-1] < model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] &&
		model.Signal.Macds[len(model.Signal.Macds)-1]-model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] >
			model.Signal.Macds[len(model.Signal.Macds)-2]-model.Signal.MacdSignals[len(model.Signal.MacdSignals)-2] &&
		!model.Signal.StochOverbought && model.Signal.StochOversold &&
		model.Signal.StochK[len(model.Signal.StochK)-1] > model.Signal.StochD[len(model.Signal.StochD)-1] {
		return true
	}

	return false
}

func CanEnterLongSlowPeriod(model *model.DataModel) bool {
	if model.Signal.Macds[len(model.Signal.Macds)-1] > model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] &&
		!model.Signal.StochOverbought && model.Signal.StochOversold &&
		model.Signal.StochK[len(model.Signal.StochK)-1] > model.Signal.StochD[len(model.Signal.StochD)-1] {
		return true
	}

	return false
}

func CanEnterShortFastPeriod(model *model.DataModel) bool {
	if model.Signal.Macds[len(model.Signal.Macds)-1] > model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] &&
		model.Signal.Macds[len(model.Signal.Macds)-1]-model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] <
			model.Signal.Macds[len(model.Signal.Macds)-2]-model.Signal.MacdSignals[len(model.Signal.MacdSignals)-2] &&
		!model.Signal.StochOversold && model.Signal.StochOverbought &&
		model.Signal.StochK[len(model.Signal.StochK)-1] < model.Signal.StochD[len(model.Signal.StochD)-1] {
		return true
	}
	return false
}

func CanEnterShortSlowPeriod(model *model.DataModel) bool {
	if model.Signal.Macds[len(model.Signal.Macds)-1] < model.Signal.MacdSignals[len(model.Signal.MacdSignals)-1] &&
		!model.Signal.StochOversold && model.Signal.StochOverbought &&
		model.Signal.StochK[len(model.Signal.StochK)-1] < model.Signal.StochD[len(model.Signal.StochD)-1] {
		return true
	}
	return false

}
