package transaction

import (
	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func UpdatePositionAfterTransaction(model *model.DataModel, order *alpaca.Order) {
	model.Position.Order = order
	model.Position.CurrentTrail = order.TrailPrice.InexactFloat64()
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

func RetrievePositionIfExists(model *model.DataModel, broker *internal.AlpacaBroker) {
	position, _ := broker.GetPosition(model.Symbol)
	if position == nil {
		model.Position.CurrentTrail = 0.0
		model.Position.HasLongPosition = false
		model.Position.HasShortPosition = false
		model.Position.FilledQuantity = 0.0
		model.Position.FilledPrice = 0.0
		model.Position.CurrentTrail = 0.0
	} else {
		order, _ := broker.RetrieveOrderIfExists(model.Symbol, "new")
		model.Position.CurrentTrail = order.TrailPrice.Abs().InexactFloat64()
		model.Position.FilledPrice = order.FilledAvgPrice.Abs().InexactFloat64()
		model.Position.FilledQuantity = order.FilledQty.Abs().InexactFloat64()
		if order.Side == alpaca.Sell {
			model.Position.HasShortPosition = true
			model.Position.HasLongPosition = false
		} else {
			model.Position.HasShortPosition = false
			model.Position.HasLongPosition = true
		}
	}
}
