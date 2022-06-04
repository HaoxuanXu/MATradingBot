package dataprocessor

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/strat/constants"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/tools"
)

func updateClose(model *model.DataModel, data *model.TotalBarData) {
	if model.CloseData.PrevMATrade == 0.0 {
		model.CloseData.PrevMATrade = data.TradeData[model.Symbol].Price
		model.CloseData.CurrMATrade = data.TradeData[model.Symbol].Price
	} else {
		model.CloseData.PrevMATrade = model.CloseData.CurrMATrade
		model.CloseData.CurrMATrade = data.TradeData[model.Symbol].Price
	}
	model.CloseData.CurrMA20Close = tools.CalcMovingAverage(data.BarData[model.Symbol], time.Now(), 20)
	model.CloseData.PrevMA20Close = tools.CalcMovingAverage(data.BarData[model.Symbol], time.Now().Add(-2*time.Minute), 20)
	model.CloseData.MASupport = tools.CalcSupportResistance(data.BarData[model.Symbol], constants.SUPPORT)
	model.CloseData.MAResistance = tools.CalcSupportResistance(data.BarData[model.Symbol], constants.RESISTANCE)
}

func ProcessBarData(model *model.DataModel, data *model.TotalBarData) {

	updateClose(model, data)
	tools.CalcAppliedTrails(model, data)

}
