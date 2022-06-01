package pipeline

import (
	"log"
	"math"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/transaction"
)

func RefreshDataModel(model *model.DataModel, broker *api.AlpacaBroker) {
	transaction.RetrievePositionIfExists(model, broker)
}

func EnterLongPosition(model *model.DataModel, broker *api.AlpacaBroker, qty float64) {
	marketOrder := broker.SubmitMarketOrder(qty, model.Symbol, "buy", "gtc")
	trail_price := math.Max(model.Trails.AppliedLongTrail, model.CloseData.CurrMAAsk*0.0015)
	trailingStopOrder := broker.SubmitTrailingStopOrder(qty, trail_price, model.Symbol, "sell")
	transaction.UpdatePositionAfterTransaction(model, marketOrder, trailingStopOrder)
	transaction.RecordEntryTransaction(model)
}

func EnterShortPosition(model *model.DataModel, broker *api.AlpacaBroker, qty float64) {
	marketOrder := broker.SubmitMarketOrder(qty, model.Symbol, "sell", "gtc")
	trail_price := math.Max(model.Trails.AppliedShortTrail, model.CloseData.CurrMABid*0.0015)
	trailingStopOrder := broker.SubmitTrailingStopOrder(qty, trail_price, model.Symbol, "buy")
	transaction.UpdatePositionAfterTransaction(model, marketOrder, trailingStopOrder)
	transaction.RecordEntryTransaction(model)
}

func ExitLongPosition(model *model.DataModel, broker *api.AlpacaBroker) {
	err := broker.ClosePosition(model.Symbol)
	if err != nil {
		log.Println(err)
	}
}

func ExitShortPosition(model *model.DataModel, broker *api.AlpacaBroker) {
	err := broker.ClosePosition(model.Symbol)
	if err != nil {
		log.Println(err)
	}
}
