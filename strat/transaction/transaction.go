package transaction

import (
	"log"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/internal/readwrite"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func UpdatePositionAfterTransaction(model *model.DataModel, order *alpaca.Order) {
	model.Position.Order = *order
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

func RetrievePositionIfExists(model *model.DataModel, broker *api.AlpacaBroker) {
	position, _ := broker.GetPosition(model.Symbol)
	if position == nil {
		model.Position.CurrentTrail = 0.0
		model.Position.HasLongPosition = false
		model.Position.HasShortPosition = false
		model.Position.HasOrder = false
	} else {
		order, _ := broker.RetrieveOrderIfExists(model.Symbol, "new")
		model.Position.Order = *order
		model.Position.HasOrder = true
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

func RecordEntryTransaction(model *model.DataModel) {
	if model.Position.HasLongPosition && !model.Position.HasShortPosition {
		log.Printf("symbol: %s, side: %s, qty: %.2f\n", model.Symbol, "buy", model.Position.FilledQuantity)
	} else if model.Position.HasShortPosition && !model.Position.HasLongPosition {
		log.Printf("symbol: %s, side: %s, qty: %.2f\n", model.Symbol, "sell", model.Position.FilledQuantity)
	}
}

func RecordExitTransaction(model *model.DataModel) {
	if !model.Position.HasShortPosition && !model.Position.HasLongPosition &&
		model.Position.HasOrder {
		log.Printf("result: $%.2f\n",
			model.Position.Order.FilledQty.Abs().InexactFloat64()*model.Position.Order.FilledAvgPrice.Abs().InexactFloat64()-
				model.Position.FilledPrice*model.Position.FilledQuantity)
	}
}

func ReadModelFromDB(model *model.DataModel) {
	model.Trails.LongTrailArray = readwrite.ReadFloatArrayToJson(model.Symbol, "long")
	model.Trails.ShortTrailArray = readwrite.ReadFloatArrayToJson(model.Symbol, "short")
}

func WriteModelToDB(model *model.DataModel) {
	readwrite.WriteFloatArrayToJson(model.Trails.LongTrailArray, model.Symbol, "long")
	readwrite.WriteFloatArrayToJson(model.Trails.ShortTrailArray, model.Symbol, "short")
}
