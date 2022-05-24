package strat

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/HaoxuanXu/MATradingBot/strat/dataprocessor"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/signalcatcher"
	"github.com/HaoxuanXu/MATradingBot/strat/transaction"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func MATradingStrategy(symbol, accountType, serverType string, entryAmount float64, totalData *model.TotalBarData) {
	broker := internal.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol, 30)
	var order *alpaca.Order

	for !broker.Clock.IsOpen {
		log.Printf("Wait for %.2f minutes till the market opens\n", time.Until(broker.Clock.NextOpen).Minutes())
	}

	for broker.Clock.IsOpen {
		dataprocessor.ProcessBarData(dataModel, totalData)
		qty := float64(int(entryAmount / dataModel.CloseData.CurrMAClose))
		if signalcatcher.CanEnterLong(dataModel) {
			order = broker.SubmitTrailingStopOrder(qty, dataModel.Position.CurrentTrail, symbol, "buy")
		}
		transaction.UpdatePositionFieldAfterTransaction(dataModel, order)
	}
}
