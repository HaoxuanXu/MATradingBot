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
	model.CloseData.CurrMA30Close = tools.CalcMovingAverage(data.Data[model.Symbol], time.Now(), 30)
	model.CloseData.PrevMA20Close = tools.CalcMovingAverage(data.Data[model.Symbol], time.Now().Add(-time.Duration(15*time.Minute)), 20)
	model.CloseData.PrevMA30Close = tools.CalcMovingAverage(data.Data[model.Symbol], time.Now().Add(-time.Duration(15*time.Minute)), 30)
}

func updateCondition(model *model.DataModel) {
	if model.CloseData.CurrMA20Close > model.CloseData.PrevMA20Close {
		model.Condition.IsMA20DaysRising = true
		model.Condition.IsMA20DaysDropping = false
	} else if model.CloseData.CurrMA20Close < model.CloseData.PrevMA20Close {
		model.Condition.IsMA20DaysRising = false
		model.Condition.IsMA20DaysDropping = true
	} else if model.CloseData.CurrMA20Close == model.CloseData.PrevMA20Close {
		model.Condition.IsMA20DaysDropping = false
		model.Condition.IsMA20DaysRising = false
	}

	if model.CloseData.CurrMA30Close > model.CloseData.PrevMA30Close {
		model.Condition.IsMA30DaysRising = true
		model.Condition.IsMA30DaysDropping = false
	} else if model.CloseData.CurrMA30Close < model.CloseData.PrevMA30Close {
		model.Condition.IsMA30DaysRising = false
		model.Condition.IsMA30DaysDropping = true
	} else if model.CloseData.CurrMA30Close == model.CloseData.PrevMA30Close {
		model.Condition.IsMA30DaysDropping = false
		model.Condition.IsMA30DaysRising = false
	}

	if model.CloseData.CurrMA20Close > model.CloseData.CurrMA30Close {
		model.Condition.IsMA20AboveMA30 = true
		model.Condition.IsMA20BelowMA30 = false
	} else if model.CloseData.CurrMA20Close < model.CloseData.CurrMA30Close {
		model.Condition.IsMA20AboveMA30 = false
		model.Condition.IsMA20BelowMA30 = true
	} else if model.CloseData.CurrMA20Close == model.CloseData.CurrMA30Close {
		model.Condition.IsMA20AboveMA30 = false
		model.Condition.IsMA20BelowMA30 = false
	}
}

func updateTrail(model *model.DataModel, data *model.TotalBarData) {
	currentBar := data.Data[model.Symbol][0]
	if model.Condition.IsMA20DaysRising {
		if currentBar.Close < model.Trails.HWM {
			model.Trails.LongTrailCandidate = math.Max(model.Trails.LongTrailCandidate, model.Trails.HWM-currentBar.Close)
		} else if currentBar.Close > model.Trails.HWM {
			model.Trails.LongTrailArray = append(model.Trails.LongTrailArray, model.Trails.LongTrailCandidate)
			model.Trails.LongTrailArray = util.ResizeFloatArray(model.Trails.LongTrailArray, model.Trails.ArrayLength)
			model.Trails.LongTrailCandidate = 0.0
			model.Trails.HWM = currentBar.Close
		}
	} else if model.Condition.IsMA20DaysDropping {
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
