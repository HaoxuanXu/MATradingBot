package datapulling

import (
	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func Pull15MinutesBarData(engine *internal.MarketDataEngine, days int) map[string][]marketdata.Bar {
	barData := engine.GetMultiBars(15, 40)
	return barData
}
