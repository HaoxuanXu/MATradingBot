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
	currentQuote := data.QuoteData[model.Symbol].AskPrice
	stop_loss := currentQuote - math.Min(math.Abs(currentQuote-model.Signal.CurrentParabolicSar), currentQuote*0.007)
	take_proft := currentQuote + math.Min(math.Abs(currentQuote-model.Signal.CurrentParabolicSar), currentQuote*0.007)*1.5
	order := broker.SubmitBracketOrder(qty, take_proft, stop_loss, model.Symbol, "buy")
	transaction.UpdatePositionAfterTransaction(model, order)
	transaction.RecordEntryTransaction(model)
}

func EnterBracketShortPosition(model *model.DataModel, data *model.TotalBarData, broker *api.AlpacaBroker, qty float64) {
	currentQuote := data.QuoteData[model.Symbol].BidPrice
	stop_loss := currentQuote + math.Min(math.Abs(currentQuote-model.Signal.CurrentParabolicSar), currentQuote*0.007)
	take_proft := currentQuote - math.Min(math.Abs(currentQuote-model.Signal.CurrentParabolicSar), currentQuote*0.007)*1.5
	order := broker.SubmitBracketOrder(qty, take_proft, stop_loss, model.Symbol, "sell")
	transaction.UpdatePositionAfterTransaction(model, order)
	transaction.RecordEntryTransaction(model)
}
