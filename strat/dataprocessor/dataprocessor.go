package dataprocessor

import (
	"math"
	"time"

	"github.com/HaoxuanXu/MATradingBot/strat/constants"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/tools"
	"github.com/HaoxuanXu/MATradingBot/util"
	"github.com/montanaflynn/stats"
)

func updateClose(model *model.DataModel, data *model.TotalBarData) {
	model.CloseData.CurrMAClose = data.Data[model.Symbol][0].Close
	model.CloseData.PrevMAClose = data.Data[model.Symbol][1].Close
	model.CloseData.CurrMA20Close = tools.CalcMovingAverage(data.Data[model.Symbol], time.Now(), 20)
	model.CloseData.PrevMA20Close = tools.CalcMovingAverage(data.Data[model.Symbol], time.Now().Add(-time.Minute), 20)
	model.CloseData.MASupport = tools.CalcSupportResistance(data.Data[model.Symbol], constants.SUPPORT)
	model.CloseData.MAResistance = tools.CalcSupportResistance(data.Data[model.Symbol], constants.RESISTANCE)
}

func updateTrail(model *model.DataModel, data *model.TotalBarData) {
	currentBar := data.Data[model.Symbol][0]

	if model.Trails.HWM == 0.0 {
		model.Trails.HWM = currentBar.High
	}

	if model.CloseData.CurrMAClose > model.CloseData.CurrMA20Close {
		model.Trails.ShortTrailCandidate = 0.0
	} else if model.CloseData.CurrMAClose < model.CloseData.CurrMA20Close {
		model.Trails.LongTrailCandidate = 0.0
	}

	if model.CloseData.CurrMAClose > model.CloseData.CurrMA20Close {
		if currentBar.High < model.Trails.HWM {
			model.Trails.LongTrailCandidate = math.Max(model.Trails.LongTrailCandidate, model.Trails.HWM-currentBar.Low)
		} else if currentBar.High > model.Trails.HWM {
			model.Trails.LongTrailArray = append(model.Trails.LongTrailArray, model.Trails.LongTrailCandidate)
			model.Trails.LongTrailArray = util.ResizeFloatArray(model.Trails.LongTrailArray, model.Trails.ArrayLength)
			model.Trails.LongTrailCandidate = 0.0
			model.Trails.HWM = currentBar.High
		}
	} else if model.CloseData.CurrMAClose < model.CloseData.CurrMA20Close {
		if currentBar.Low > model.Trails.HWM {
			model.Trails.ShortTrailCandidate = math.Max(model.Trails.ShortTrailCandidate, currentBar.High-model.Trails.HWM)
		} else if currentBar.Low < model.Trails.HWM {
			model.Trails.ShortTrailArray = append(model.Trails.ShortTrailArray, model.Trails.ShortTrailCandidate)
			model.Trails.ShortTrailArray = util.ResizeFloatArray(model.Trails.ShortTrailArray, model.Trails.ArrayLength)
			model.Trails.ShortTrailCandidate = 0.0
			model.Trails.HWM = currentBar.Low
		}
	}

	if len(model.Trails.LongTrailArray) >= model.Trails.ArrayLength {
		model.Trails.AppliedLongTrail, _ = stats.Percentile(model.Trails.LongTrailArray, 95.0)
	}
	if len(model.Trails.ShortTrailArray) >= model.Trails.ArrayLength {
		model.Trails.AppliedShortTrail, _ = stats.Percentile(model.Trails.ShortTrailArray, 95.0)
	}
}

func ProcessBarData(model *model.DataModel, data *model.TotalBarData) {

	updateClose(model, data)
	updateTrail(model, data)

}
