package dataprocessor

import (
	"math"
	"time"

	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/tools"
	"github.com/HaoxuanXu/MATradingBot/util"
	"github.com/montanaflynn/stats"
)

func fillCurrPrevClose(model *model.DataModel, data *model.TotalBarData) {
	model.CloseData.CurrMAClose = data.Data[model.Symbol][0].Close
	model.CloseData.CurrMA20Close = tools.CalcMovingAverage(data.Data[model.Symbol], time.Now(), 20)
	model.CloseData.PrevMA20Close = tools.CalcMovingAverage(data.Data[model.Symbol], time.Now().Add(-time.Duration(15*time.Minute)), 20)
}

func updateCondition(model *model.DataModel) {
	model.Condition.IsMA20PeriodsPreviouslyDropping = model.Condition.IsMA20PeriodsDropping
	model.Condition.IsMA20PeriodsPreviouslyRising = model.Condition.IsMA20PeriodsRising

	if model.CloseData.CurrMA20Close > model.CloseData.PrevMA20Close {
		model.Condition.IsMA20PeriodsRising = true
		model.Condition.IsMA20PeriodsDropping = false
	} else if model.CloseData.CurrMA20Close < model.CloseData.PrevMA20Close {
		model.Condition.IsMA20PeriodsRising = false
		model.Condition.IsMA20PeriodsDropping = true
	} else if model.CloseData.CurrMA20Close == model.CloseData.PrevMA20Close {
		model.Condition.IsMA20PeriodsDropping = false
		model.Condition.IsMA20PeriodsRising = false
	}

	if model.CloseData.CurrMAClose > model.CloseData.CurrMA20Close {
		model.Condition.IsMAAboveMA20 = true
		model.Condition.IsMABelowMA20 = false
	} else {
		model.Condition.IsMABelowMA20 = true
		model.Condition.IsMAAboveMA20 = false
	}

}

func updateTrail(model *model.DataModel, data *model.TotalBarData) {
	currentBar := data.Data[model.Symbol][0]

	if model.Condition.IsMA20PeriodsRising && !model.Condition.IsMA20PeriodsPreviouslyRising {
		model.Trails.ShortTrailCandidate = 0.0
	} else if model.Condition.IsMA20PeriodsDropping && !model.Condition.IsMA20PeriodsPreviouslyDropping {
		model.Trails.LongTrailCandidate = 0.0
	}

	if model.Condition.IsMA20PeriodsRising {
		if currentBar.Close < model.Trails.HWM {
			model.Trails.LongTrailCandidate = math.Max(model.Trails.LongTrailCandidate, model.Trails.HWM-currentBar.Close)
		} else if currentBar.Close > model.Trails.HWM {
			model.Trails.LongTrailArray = append(model.Trails.LongTrailArray, model.Trails.LongTrailCandidate)
			model.Trails.LongTrailArray = util.ResizeFloatArray(model.Trails.LongTrailArray, model.Trails.ArrayLength)
			model.Trails.LongTrailCandidate = 0.0
			model.Trails.HWM = currentBar.Close
		}
	} else if model.Condition.IsMA20PeriodsDropping {
		if currentBar.Close > model.Trails.HWM {
			model.Trails.ShortTrailCandidate = math.Max(model.Trails.ShortTrailCandidate, currentBar.Close-model.Trails.HWM)
		} else if currentBar.Close < model.Trails.HWM {
			model.Trails.ShortTrailArray = append(model.Trails.ShortTrailArray, model.Trails.ShortTrailCandidate)
			model.Trails.ShortTrailArray = util.ResizeFloatArray(model.Trails.ShortTrailArray, model.Trails.ArrayLength)
			model.Trails.ShortTrailCandidate = 0.0
			model.Trails.HWM = currentBar.Close
		}
	}
	model.Trails.AppliedLongTrail, _ = stats.Percentile(model.Trails.LongTrailArray, 95.0)
	model.Trails.AppliedShortTrail, _ = stats.Percentile(model.Trails.ShortTrailArray, 95.0)
}

func ProcessBarData(model *model.DataModel, data *model.TotalBarData) {
	fillCurrPrevClose(model, data)
	updateCondition(model)
	updateTrail(model, data)

}
