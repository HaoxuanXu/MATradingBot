package dataprocessor

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/tools"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func fillCurrPrevClose(symbol string, barResult map[string][]marketdata.Bar, barData *model.MABarCloseData) {
	barData.CurrMA20Close = tools.CalcMovingAverage(barResult[symbol], time.Now(), 20)
	barData.CurrMA30Close = tools.CalcMovingAverage(barResult[symbol], time.Now(), 30)
	barData.PrevMA20Close = tools.CalcMovingAverage(barResult[symbol], time.Now().Add(-time.Duration(15*time.Minute)), 20)
	barData.PrevMA30Close = tools.CalcMovingAverage(barResult[symbol], time.Now().Add(-time.Duration(15*time.Minute)), 30)
}

func updateCondition(barData *model.MABarCloseData, condition *model.MAConditions) {
	if barData.CurrMA20Close > barData.PrevMA20Close {
		condition.IsMA20DaysRising = true
		condition.IsMA20DaysDropping = false
	} else if barData.CurrMA20Close < barData.PrevMA20Close {
		condition.IsMA20DaysRising = false
		condition.IsMA20DaysDropping = true
	} else if barData.CurrMA20Close == barData.PrevMA20Close {
		condition.IsMA20DaysDropping = false
		condition.IsMA20DaysRising = false
	}

	if barData.CurrMA30Close > barData.PrevMA30Close {
		condition.IsMA30DaysRising = true
		condition.IsMA30DaysDropping = false
	} else if barData.CurrMA30Close < barData.PrevMA30Close {
		condition.IsMA30DaysRising = false
		condition.IsMA30DaysDropping = true
	} else if barData.CurrMA30Close == barData.PrevMA30Close {
		condition.IsMA30DaysDropping = false
		condition.IsMA30DaysRising = false
	}

	if barData.CurrMA20Close > barData.CurrMA30Close {
		condition.IsMA20AboveMA30 = true
		condition.IsMA20BelowMA30 = false
	} else if barData.CurrMA20Close < barData.CurrMA30Close {
		condition.IsMA20AboveMA30 = false
		condition.IsMA20BelowMA30 = true
	} else if barData.CurrMA20Close == barData.CurrMA30Close {
		condition.IsMA20AboveMA30 = false
		condition.IsMA20BelowMA30 = false
	}
}

func ProcessBarData(symbol string, barResult map[string][]marketdata.Bar, barData *model.MABarCloseData, condition *model.MAConditions) {
	fillCurrPrevClose(symbol, barResult, barData)
	updateCondition(barData, condition)

}
