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
	model.CloseData.CurrMAAsk = data.QuoteData[model.Symbol].AskPrice
	model.CloseData.CurrMABid = data.QuoteData[model.Symbol].BidPrice
	model.CloseData.CurrMA20Close = tools.CalcMovingAverage(data.BarData[model.Symbol], time.Now(), 20)
	model.CloseData.PrevMA20Close = tools.CalcMovingAverage(data.BarData[model.Symbol], time.Now().Add(-2*time.Minute), 20)
	model.CloseData.MASupport = tools.CalcSupportResistance(data.BarData[model.Symbol], constants.SUPPORT)
	model.CloseData.MAResistance = tools.CalcSupportResistance(data.BarData[model.Symbol], constants.RESISTANCE)
}

func updateTrail(model *model.DataModel, data *model.TotalBarData) {
	currentBar := data.BarData[model.Symbol][0]
	if model.Trails.LongHWM == 0.0 || model.Trails.ShortHWM == 0.0 {
		model.Trails.LongHWM = currentBar.High
		model.Trails.ShortHWM = currentBar.Low
	}

	if currentBar.Close > model.CloseData.CurrMA20Close {
		model.Trails.ShortTrailCandidate = 0.0
		model.Trails.ShortHWM = currentBar.Low
	} else if currentBar.Close < model.CloseData.CurrMA20Close {
		model.Trails.LongTrailCandidate = 0.0
		model.Trails.LongHWM = currentBar.High
	}

	if currentBar.Low > model.CloseData.CurrMA20Close {
		if currentBar.High < model.Trails.LongHWM {
			model.Trails.LongTrailCandidate = math.Max(model.Trails.LongTrailCandidate, model.Trails.LongHWM-currentBar.Low)
		} else if currentBar.High > model.Trails.LongHWM {
			if model.Trails.LongTrailCandidate > 0 {
				model.Trails.LongTrailArray = append(model.Trails.LongTrailArray, model.Trails.LongTrailCandidate)
				model.Trails.LongTrailArray = util.ResizeFloatArray(model.Trails.LongTrailArray, model.Trails.ArrayLength)
			}
			model.Trails.LongTrailCandidate = 0.0
			model.Trails.LongHWM = currentBar.High
		}
	} else if currentBar.High < model.CloseData.CurrMA20Close {
		if currentBar.Low > model.Trails.ShortHWM {
			model.Trails.ShortTrailCandidate = math.Max(model.Trails.ShortTrailCandidate, currentBar.High-model.Trails.ShortHWM)
		} else if currentBar.Low < model.Trails.ShortHWM {
			if model.Trails.ShortTrailCandidate > 0 {
				model.Trails.ShortTrailArray = append(model.Trails.ShortTrailArray, model.Trails.ShortTrailCandidate)
				model.Trails.ShortTrailArray = util.ResizeFloatArray(model.Trails.ShortTrailArray, model.Trails.ArrayLength)
			}
			model.Trails.ShortTrailCandidate = 0.0
			model.Trails.ShortHWM = currentBar.Low
		}
	}
	// log.Printf("%s long hwm: %.2f; short hwm: %.2f; current high: %.2f; current low: %.2f; long trail: %.2f; short trail: %.2f; timestamp: %s\n",
	// 	model.Symbol, model.Trails.LongHWM, model.Trails.ShortHWM, currentBar.High,
	// 	currentBar.Low, model.Trails.LongTrailCandidate, model.Trails.ShortTrailCandidate, currentBar.Timestamp.String())

	if len(model.Trails.LongTrailArray) > 0 {
		model.Trails.AppliedLongTrail, _ = stats.Max(model.Trails.LongTrailArray)
	}
	if len(model.Trails.ShortTrailArray) >= 0 {
		model.Trails.AppliedShortTrail, _ = stats.Max(model.Trails.ShortTrailArray)
	}

}

func ProcessBarData(model *model.DataModel, data *model.TotalBarData) {

	updateClose(model, data)
	updateTrail(model, data)

}
