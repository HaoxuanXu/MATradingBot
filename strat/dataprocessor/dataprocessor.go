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
		model.Condition.HWM = math.Max(currentBar.Close, model.Condition.HWM)
		if currentBar.Close < model.Condition.HWM {
			model.Trails.TrailDataArray = append(model.Trails.TrailDataArray, model.Condition.HWM-currentBar.Close)
			model.Trails.TrailDataArray = util.ResizeFloatArray(model.Trails.TrailDataArray, model.Trails.DataLength)
		}
	} else if model.Condition.IsMA20DaysDropping {
		model.Condition.HWM = math.Min(currentBar.Close, model.Condition.HWM)
		if currentBar.Close > model.Condition.HWM {
			model.Trails.TrailDataArray = append(model.Trails.TrailDataArray, currentBar.Close-model.Condition.HWM)
			model.Trails.TrailDataArray = util.ResizeFloatArray(model.Trails.TrailDataArray, model.Trails.DataLength)
		}
	}
	model.Position.CurrentTrail, _ = stats.Max(model.Trails.TrailDataArray)
}

func ProcessBarData(model *model.DataModel, data *model.TotalBarData) {
	fillCurrPrevClose(model, data)
	updateCondition(model)
	updateTrail(model, data)

}
