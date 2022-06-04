package transaction

import (
	"log"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func UpdatePositionAfterTransaction(model *model.DataModel, marketOrder, trailingStopOrder *alpaca.Order) {
	model.Position.MarketOrder = *marketOrder
	model.Position.TrailingStopOrder = *trailingStopOrder
	model.Position.CurrentTrail = trailingStopOrder.TrailPrice.InexactFloat64()
	model.Position.FilledQuantity = marketOrder.FilledQty.Abs().InexactFloat64()
	model.Position.FilledPrice = marketOrder.FilledAvgPrice.InexactFloat64()
	if marketOrder.Side == alpaca.Sell {
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
		model.Position.CurrentTrail = 0.0
		model.Position.HasLongPosition = false
		model.Position.HasShortPosition = false
		model.Position.HasOrder = false
	} else {
		marketOrder, err := broker.RetrieveOrderIfExists(model.Symbol, "filled", "market")
		if err != nil {
			log.Println(err)
		}
		trailingStopOrder, err := broker.RetrieveOrderIfExists(model.Symbol, "open", "trailing_stop")
		if err != nil {
			log.Println(err)
		}
		if marketOrder != nil {
			model.Position.HasOrder = true
			model.Position.MarketOrder = *marketOrder
			model.Position.FilledPrice = marketOrder.FilledAvgPrice.Abs().InexactFloat64()
			model.Position.FilledQuantity = marketOrder.FilledQty.Abs().InexactFloat64()
			if marketOrder.Side == alpaca.Sell {
				model.Position.HasShortPosition = true
				model.Position.HasLongPosition = false
			} else {
				model.Position.HasShortPosition = false
				model.Position.HasLongPosition = true
			}
		}
		if trailingStopOrder != nil {
			model.Position.TrailingStopOrder = *trailingStopOrder
			model.Position.CurrentTrail = trailingStopOrder.TrailPrice.Abs().InexactFloat64()
		}
	}
}

func RecordEntryTransaction(model *model.DataModel) {
	if model.Position.HasLongPosition && !model.Position.HasShortPosition {
		log.Printf("symbol: %s, side: %s, qty: %.2f\n", model.Symbol, "buy", model.Position.FilledQuantity)
	} else if model.Position.HasShortPosition && !model.Position.HasLongPosition {
		log.Printf("symbol: %s, side: %s, qty: %.2f\n", model.Symbol, "sell", model.Position.FilledQuantity)
	}
}

func RecordExitTransaction(model *model.DataModel) {
	if !model.Position.HasShortPosition && !model.Position.HasLongPosition {
		log.Printf("result: $%.2f\n",
			model.Position.MarketOrder.FilledQty.Abs().InexactFloat64()*model.Position.MarketOrder.FilledAvgPrice.Abs().InexactFloat64()-
				model.Position.FilledPrice*model.Position.FilledQuantity)
	}
}
