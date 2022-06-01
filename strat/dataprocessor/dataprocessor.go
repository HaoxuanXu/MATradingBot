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
	model.CloseData.MASupport = tools.CalcSupportResistance(data.BarData[model.Symbol], constants.SUPPORT)
	model.CloseData.MAResistance = tools.CalcSupportResistance(data.BarData[model.Symbol], constants.RESISTANCE)
}

func updateTrail(model *model.DataModel, data *model.TotalBarData) {
	currentBar := data.BarData[model.Symbol][0]
	if model.Trails.LongHWM == 0.0 || model.Trails.ShortHWM == 0.0 {
		model.Trails.LongHWM = currentBar.High
		model.Trails.ShortHWM = currentBar.Low
	}

	if model.CloseData.CurrMAAsk > model.CloseData.CurrMA20Close {
		model.Trails.ShortTrailCandidate = 0.0
	} else if model.CloseData.CurrMABid < model.CloseData.CurrMA20Close {
		model.Trails.LongTrailCandidate = 0.0
	}

	if model.CloseData.CurrMAAsk > model.CloseData.CurrMA20Close {
		if model.CloseData.CurrMAAsk < model.Trails.LongHWM {
			model.Trails.LongTrailCandidate = math.Max(model.Trails.LongTrailCandidate, model.Trails.LongHWM-currentBar.Low)
		} else if model.CloseData.CurrMAAsk > model.Trails.LongHWM {
			if model.Trails.LongTrailCandidate > 0 {
				model.Trails.LongTrailArray = append(model.Trails.LongTrailArray, model.Trails.LongTrailCandidate)
				model.Trails.LongTrailArray = util.ResizeFloatArray(model.Trails.LongTrailArray, model.Trails.ArrayLength)
			}
			model.Trails.LongTrailCandidate = 0.0
			model.Trails.LongHWM = model.CloseData.CurrMAAsk
		}
	} else if model.CloseData.CurrMABid < model.CloseData.CurrMA20Close {
		if model.CloseData.CurrMABid > model.Trails.ShortHWM {
			model.Trails.ShortTrailCandidate = math.Max(model.Trails.ShortTrailCandidate, currentBar.High-model.Trails.ShortHWM)
		} else if model.CloseData.CurrMABid < model.Trails.ShortHWM {
			if model.Trails.ShortTrailCandidate > 0 {
				model.Trails.ShortTrailArray = append(model.Trails.ShortTrailArray, model.Trails.ShortTrailCandidate)
				model.Trails.ShortTrailArray = util.ResizeFloatArray(model.Trails.ShortTrailArray, model.Trails.ArrayLength)
			}
			model.Trails.ShortTrailCandidate = 0.0
			model.Trails.ShortHWM = model.CloseData.CurrMABid
		}
	}
	// log.Printf("%s long hwm: %.2f; short hwm: %.2f; current high: %.2f; current low: %.2f; long trail: %.2f; short trail: %.2f; timestamp: %s\n",
	// 	model.Symbol, model.Trails.LongHWM, model.Trails.ShortHWM, currentBar.High,
	// 	currentBar.Low, model.Trails.LongTrailCandidate, model.Trails.ShortTrailCandidate, currentBar.Timestamp.String())

	if len(model.Trails.LongTrailArray) >= 5 {
		model.Trails.AppliedLongTrail, _ = stats.Percentile(model.Trails.LongTrailArray, 95.0)
	}
	if len(model.Trails.ShortTrailArray) >= 5 {
		model.Trails.AppliedShortTrail, _ = stats.Percentile(model.Trails.ShortTrailArray, 95.0)
	}

}

func ProcessBarData(model *model.DataModel, data *model.TotalBarData) {

	updateClose(model, data)
	updateTrail(model, data)

}
