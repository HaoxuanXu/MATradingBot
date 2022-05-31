package strat

import (
	"log"
	"math/rand"
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strat/dataprocessor"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/pipeline"
	"github.com/HaoxuanXu/MATradingBot/strat/signalcatcher"
	"github.com/HaoxuanXu/MATradingBot/strat/transaction"
)

func MATradingStrategy(symbol, accountType, serverType string, entryPercent float64, totalData *model.TotalBarData, channel chan bool) {
	broker := api.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol, 20)
	transaction.ReadModelFromDB(dataModel)
	entryAmount := broker.Cash * entryPercent

	for <-channel {
		// Sleep for a random seconds between 0 and 10 seconds
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(10)) * time.Minute)

		dataprocessor.ProcessBarData(dataModel, totalData)
		pipeline.RefreshDataModel(dataModel, broker)
		transaction.RecordExitTransaction(dataModel)

		qty := float64(int(entryAmount / dataModel.CloseData.CurrMAClose))
		if signalcatcher.CanEnterLong(dataModel, broker) {
			pipeline.EnterLongPosition(dataModel, broker, qty)
		} else if signalcatcher.CanEnterShort(dataModel, broker) {
			pipeline.EnterShortPosition(dataModel, broker, qty)
		}
		transaction.WriteModelToDB(dataModel)
	}
	log.Printf("%s worker closed", symbol)

}
