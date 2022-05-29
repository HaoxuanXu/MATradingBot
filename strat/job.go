package strat

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/HaoxuanXu/MATradingBot/strat/dataprocessor"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/pipeline"
	"github.com/HaoxuanXu/MATradingBot/strat/signalcatcher"
	"github.com/HaoxuanXu/MATradingBot/strat/transaction"
)

func MATradingStrategy(symbol, accountType, serverType string, entryAmount float64, totalData *model.TotalBarData) {
	broker := internal.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol, 30)
	transaction.RetrievePositionIfExists(dataModel, broker)
	transaction.ReadModelFromDB(dataModel)

	for !broker.Clock.IsOpen {
		log.Printf("Wait for %.2f minutes till the market opens\n", time.Until(broker.Clock.NextOpen).Minutes())
	}

	for broker.Clock.IsOpen {
		dataprocessor.ProcessBarData(dataModel, totalData)
		qty := float64(int(entryAmount / dataModel.CloseData.CurrMAClose))
		if signalcatcher.CanEnterLong(dataModel, broker) {
			pipeline.EnterLongPosition(dataModel, broker, qty)
		} else if signalcatcher.CanEnterShort(dataModel, broker) {
			pipeline.EnterShortPosition(dataModel, broker, qty)
		} else {
			time.Sleep(time.Minute)
			transaction.WriteModelToDB(dataModel)
		}
	}

}
