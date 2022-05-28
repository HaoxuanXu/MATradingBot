package pipeline

import (
	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/transaction"
)

func PopulateDataModel(model *model.DataModel, broker *internal.AlpacaBroker) {
	transaction.RetrievePositionIfExists(model, broker)
}

func EnterLongPosition(model *model.DataModel, broker *internal.AlpacaBroker, qty float64) {
	order := broker.SubmitTrailingStopOrder(qty, model.Trails.AppliedLongTrail, model.Symbol, "buy")
	transaction.UpdatePositionAfterTransaction(model, order)
}

func EnterShortPosition(model *model.DataModel, broker *internal.AlpacaBroker, qty float64) {
	order := broker.SubmitTrailingStopOrder(qty, model.Trails.AppliedShortTrail, model.Symbol, "sell")
	transaction.UpdatePositionAfterTransaction(model, order)
}
