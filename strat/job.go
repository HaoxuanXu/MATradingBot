package strat

import (
	"log"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strat/dataprocessor"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/pipeline"
	"github.com/HaoxuanXu/MATradingBot/strat/signalcatcher"
	"github.com/HaoxuanXu/MATradingBot/strat/transaction"
	"github.com/HaoxuanXu/MATradingBot/util"
)

func MATradingStrategy(symbol, accountType, serverType string, entryPercent float64, totalData *model.TotalBarData, channel chan bool) {
	defer util.HandlePanic()
	broker := api.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol, 100)
	transaction.ReadModelFromDB(dataModel)
	entryAmount := broker.Cash * entryPercent

	for <-channel {

		dataprocessor.ProcessBarData(dataModel, totalData)

		pipeline.RefreshDataModel(dataModel, broker)
		transaction.WriteModelToDB(dataModel)

		qty := float64(int(entryAmount / dataModel.CloseData.CurrMAAsk))

		if signalcatcher.CanEnterLong(dataModel, broker) {
			pipeline.EnterLongPosition(dataModel, broker, qty)
		} else if signalcatcher.CanEnterShort(dataModel, broker) {
			pipeline.EnterShortPosition(dataModel, broker, qty)
		} else if signalcatcher.CanExitLong(dataModel, broker) {
			pipeline.ExitLongPosition(dataModel, broker)
		} else if signalcatcher.CanExitShort(dataModel, broker) {
			pipeline.ExitShortPosition(dataModel, broker)
		}

	}
	log.Printf("%s worker closed", symbol)

}
