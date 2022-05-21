package dataprocessor

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/tools"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func fillCurrPrevClose(symbol string, barResult map[string][]marketdata.Bar, barData *model.MABarData) {
	barData.CurrMA20Close = tools.CalcMovingAverage(barResult[symbol], time.Now(), 20)
	barData.CurrMA30Close = tools.CalcMovingAverage(barResult[symbol], time.Now(), 30)
}
