package pipeline

import (
	"math"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/transaction"
)

func RefreshPosition(model *model.DataModel, broker *api.AlpacaBroker) {
	transaction.RetrievePositionIfExists(model, broker)
}

func EnterBracketLongPosition(model *model.DataModel, data *model.TotalBarData, broker *api.AlpacaBroker, qty float64) {
	currentTrade := data.TradeData[model.Symbol].Price
	stop_loss := currentTrade - math.Abs(currentTrade-model.Signal.CurrentParabolicSar)
	take_proft := currentTrade + math.Abs(currentTrade-model.Signal.CurrentParabolicSar)
	order := broker.SubmitBracketOrder(qty, take_proft, stop_loss, model.Symbol, "buy")
	transaction.UpdatePositionAfterTransaction(model, order)
	transaction.RecordEntryTransaction(model)
}

func EnterBracketShortPosition(model *model.DataModel, data *model.TotalBarData, broker *api.AlpacaBroker, qty float64) {
	currentTrade := data.TradeData[model.Symbol].Price
	stop_loss := currentTrade + math.Abs(currentTrade-model.Signal.CurrentParabolicSar)
	take_proft := currentTrade - math.Abs(currentTrade-model.Signal.CurrentParabolicSar)
	order := broker.SubmitBracketOrder(qty, take_proft, stop_loss, model.Symbol, "sell")
	transaction.UpdatePositionAfterTransaction(model, order)
	transaction.RecordEntryTransaction(model)
}
