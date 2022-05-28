package datapulling

import (
	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func PullMinuteBarData(engine *internal.MarketDataEngine, days int) map[string][]marketdata.Bar {
	barData := engine.GetMultiBars(1)
	return barData
}
