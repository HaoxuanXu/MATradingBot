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
	order := broker.SubmitTrailingStopOrder(qty, model.Trails.AppliedLongTrail, model.Symbol, "buy")
	transaction.UpdatePositionAfterTransaction(model, order)
	transaction.RecordEntryTransaction(model)
}

func EnterShortPosition(model *model.DataModel, broker *api.AlpacaBroker, qty float64) {
	order := broker.SubmitTrailingStopOrder(qty, model.Trails.AppliedShortTrail, model.Symbol, "sell")
	transaction.UpdatePositionAfterTransaction(model, order)
	transaction.RecordEntryTransaction(model)
}
