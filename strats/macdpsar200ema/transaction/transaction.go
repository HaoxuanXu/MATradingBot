package transaction

import (
	"log"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func UpdatePositionAfterTransaction(model *model.DataModel, order *alpaca.Order) {
	model.Position.Order = *order
	model.Position.FilledQuantity = order.FilledQty.Abs().InexactFloat64()
	model.Position.FilledPrice = order.FilledAvgPrice.InexactFloat64()
	if order.Side == alpaca.Sell {
		model.Position.HasShortPosition = true
		model.Position.HasLongPosition = false
	} else {
		model.Position.HasLongPosition = true
		model.Position.HasShortPosition = false
	}
}

func RetrievePositionIfExists(model *model.DataModel, broker *api.AlpacaBroker) {
	position, _ := broker.GetPosition(model.Symbol)

	if position == nil {
		broker.HasLongPosition = false
		broker.HasShortPosition = false
		model.Position.HasOrder = false
	} else {
		if position.Side == "sell" {
			broker.HasShortPosition = true
			broker.HasLongPosition = false
		} else {
			broker.HasShortPosition = false
			broker.HasLongPosition = true
		}
		marketOrder, err := broker.RetrieveOrderIfExists(model.Symbol, "filled", "market")
		if err != nil {
			log.Println(err)
		}
		if err != nil {
			log.Println(err)
		}
		if marketOrder != nil {
			model.Position.HasOrder = true
			model.Position.Order = *marketOrder
			model.Position.FilledPrice = marketOrder.FilledAvgPrice.Abs().InexactFloat64()
			model.Position.FilledQuantity = marketOrder.FilledQty.Abs().InexactFloat64()
		}
	}
}

func RecordEntryTransaction(model *model.DataModel, broker *api.AlpacaBroker) {
	if broker.HasLongPosition && !broker.HasShortPosition {
		log.Printf("symbol: %s, side: %s, qty: %.2f\n", model.Symbol, "buy", model.Position.FilledQuantity)
	} else if broker.HasShortPosition && !broker.HasLongPosition {
		log.Printf("symbol: %s, side: %s, qty: %.2f\n", model.Symbol, "sell", model.Position.FilledQuantity)
	}
}

func RecordExitTransaction(model *model.DataModel) {
	if !model.Position.HasShortPosition && !model.Position.HasLongPosition {
		log.Printf("result: $%.2f\n",
			model.Position.Order.FilledQty.Abs().InexactFloat64()*model.Position.Order.FilledAvgPrice.Abs().InexactFloat64()-
				model.Position.FilledPrice*model.Position.FilledQuantity)
	}
}
