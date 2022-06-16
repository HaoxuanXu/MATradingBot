package pipeline

import (
	"log"
	"math"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/transaction"
)

func RefreshPosition(model *model.DataModel, broker *api.AlpacaBroker) {
	transaction.RetrievePositionIfExists(model, broker)
}

func EnterTrailingStopLongPosition(model *model.DataModel, broker *api.AlpacaBroker, qty float64) {
	quote := model.Signal.Quote.AskPrice
	trail := math.Max(quote*0.001, math.Abs(model.Signal.TrailingStopLossLong-quote)/2)
	marketOrder, err := broker.SubmitMarketOrder(
		qty,
		model.Symbol,
		"buy",
		"day",
	)
	if err != nil {
		log.Printf("%s (market order): %v\n", model.Symbol, err)
	}
	_, err = broker.SubmitTrailingStopOrder(
		qty,
		trail,
		model.Symbol,
		"sell",
	)
	if err != nil {
		log.Printf("%s (trailing order): %v\n", model.Symbol, err)
		return
	}

	transaction.UpdatePositionAfterTransaction(model, marketOrder, broker)
	transaction.RecordEntryTransaction(model, broker)
}

func EnterTrailingStopShortPosition(model *model.DataModel, broker *api.AlpacaBroker, qty float64) {
	quote := model.Signal.Quote.BidPrice
	trail := math.Max(quote*0.0012, math.Abs(model.Signal.TrailingStopLossShort-quote))
	marketOrder, err := broker.SubmitMarketOrder(
		qty,
		model.Symbol,
		"sell",
		"day",
	)
	if err != nil {
		log.Printf("%s (market order): %v\n", model.Symbol, err)
	}
	_, err = broker.SubmitTrailingStopOrder(
		qty,
		trail,
		model.Symbol,
		"buy",
	)
	if err != nil {
		log.Printf("%s (trailing order): %v\n", model.Symbol, err)
		return
	}

	transaction.UpdatePositionAfterTransaction(model, marketOrder, broker)
	transaction.RecordEntryTransaction(model, broker)
}

func EnterBracketLongPosition(model *model.DataModel, data *model.TotalBarData, broker *api.AlpacaBroker, qty float64) {
	currentQuote := data.StockQuoteData[model.Symbol].AskPrice
	profitOffset := math.Min(math.Abs(model.Signal.Bars[len(model.Signal.Bars)-1].Low-model.Signal.ParabolicSars[len(model.Signal.ParabolicSars)-1]),
		currentQuote*0.005)
	if profitOffset < 0.01 {
		log.Printf("long %s: profit offset value is only $%f, lower than $0.01 minimum\n", model.Symbol, profitOffset)
		return
	}
	stop_loss := currentQuote - profitOffset
	take_profit := currentQuote + profitOffset
	if take_profit > stop_loss {
		order, err := broker.SubmitBracketOrder(qty, take_profit, stop_loss, model.Symbol, "buy")
		if err != nil {
			log.Printf("%s: %v\n", model.Symbol, err)
			return
		}
		transaction.UpdatePositionAfterTransaction(model, order, broker)
		transaction.RecordEntryTransaction(model, broker)
	} else {
		log.Printf("long %s: take profit only $%.2f while stop loss is $%.2f; take profit smaller than stop loss\n",
			model.Symbol, take_profit, stop_loss)
	}

}

func EnterBracketShortPosition(model *model.DataModel, data *model.TotalBarData, broker *api.AlpacaBroker, qty float64) {

	currentQuote := data.StockQuoteData[model.Symbol].BidPrice
	profitOffset := math.Abs(math.Min(math.Abs(model.Signal.ParabolicSars[len(model.Signal.ParabolicSars)-1]-model.Signal.Bars[len(model.Signal.Bars)-1].High),
		currentQuote*0.005))
	if profitOffset < 0.01 {
		log.Printf("short %s: profit offset value is only $%f, lower than $0.01 minimum\n", model.Symbol, profitOffset)
		return
	}
	stop_loss := currentQuote + profitOffset
	take_profit := currentQuote - profitOffset
	if take_profit < stop_loss {
		order, err := broker.SubmitBracketOrder(qty, take_profit, stop_loss, model.Symbol, "sell")
		if err != nil {
			log.Printf("%s: %v\n", model.Symbol, err)
			return
		}
		transaction.UpdatePositionAfterTransaction(model, order, broker)
		transaction.RecordEntryTransaction(model, broker)
	} else {
		log.Printf("short %s: take profit only $%.2f while stop loss is $%.2f; take profit larger than stop loss\n",
			model.Symbol, take_profit, stop_loss)
	}
}
