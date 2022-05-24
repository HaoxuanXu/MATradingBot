package transaction

import (
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func UpdatePositionFieldAfterTransaction(model *model.DataModel, order *alpaca.Order) {
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
