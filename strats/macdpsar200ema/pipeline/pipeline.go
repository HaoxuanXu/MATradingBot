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

func EnterBracketLongPosition(model *model.DataModel, data *model.TotalBarData, broker *api.AlpacaBroker, qty float64, isCrypto bool) {
	var profitOffset float64
	var currentQuote float64
	if !isCrypto {
		currentQuote = data.StockQuoteData[model.Symbol].AskPrice
		profitOffset = math.Min(math.Abs(model.Signal.CurrentBar.Low-model.Signal.CurrentParabolicSar),
			currentQuote*0.007)
		if profitOffset < 0.01 {
			return
		}
	} else {
		currentQuote = data.CryptoQuoteData[model.Symbol].AskPrice
		profitOffset = math.Abs(math.Min(math.Abs(model.Signal.CurrentCryptoBar.Low-model.Signal.CurrentParabolicSar),
			currentQuote*0.007))

	}

	stop_loss := currentQuote - profitOffset
	take_profit := currentQuote + profitOffset
	if take_profit > stop_loss {
		order := broker.SubmitBracketOrder(qty, take_profit, stop_loss, model.Symbol, "buy")
		transaction.UpdatePositionAfterTransaction(model, order)
		transaction.RecordEntryTransaction(model)
	}

}

func EnterBracketShortPosition(model *model.DataModel, data *model.TotalBarData, broker *api.AlpacaBroker, qty float64, isCrypto bool) {
	var currentQuote float64
	var profitOffset float64
	if !isCrypto {
		currentQuote = data.StockQuoteData[model.Symbol].BidPrice
		profitOffset = math.Abs(math.Min(math.Abs(model.Signal.CurrentParabolicSar-model.Signal.CurrentBar.High), currentQuote*0.007))
		if profitOffset < 0.01 {
			return
		}
		stop_loss := currentQuote + profitOffset
		take_profit := currentQuote - profitOffset
		if take_profit < stop_loss {
			order := broker.SubmitBracketOrder(qty, take_profit, stop_loss, model.Symbol, "sell")
			transaction.UpdatePositionAfterTransaction(model, order)
			transaction.RecordEntryTransaction(model)
		}
	} else {
		// Crypto API current does not accept short
		return
	}
}
