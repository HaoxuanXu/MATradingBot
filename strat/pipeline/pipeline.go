package pipeline

import (
	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/transaction"
)

func RefreshDataModel(model *model.DataModel, broker *api.AlpacaBroker) {
	transaction.RetrievePositionIfExists(model, broker)
}

func EnterLongPosition(model *model.DataModel, broker *api.AlpacaBroker, qty float64) {
	marketOrder := broker.SubmitMarketOrder(qty, model.Symbol, "buy", "gtc")
	trailingStopOrder := broker.SubmitTrailingStopOrder(qty, model.Trails.AppliedLongTrail, model.Symbol, "sell")
	transaction.UpdatePositionAfterTransaction(model, marketOrder, trailingStopOrder)
	transaction.RecordEntryTransaction(model)
}

func EnterShortPosition(model *model.DataModel, broker *api.AlpacaBroker, qty float64) {
	marketOrder := broker.SubmitMarketOrder(qty, model.Symbol, "sell", "gtc")
	trailingStopOrder := broker.SubmitTrailingStopOrder(qty, model.Trails.AppliedShortTrail, model.Symbol, "buy")
	transaction.UpdatePositionAfterTransaction(model, marketOrder, trailingStopOrder)
	transaction.RecordEntryTransaction(model)
}
