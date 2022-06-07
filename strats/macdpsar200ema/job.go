package macdpsar200ema

import (
	"log"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/dataprocessor"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/pipeline"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/signalcatcher"
	"github.com/HaoxuanXu/MATradingBot/util"
)

func MACDPSar200EMAStrategy(symbol, accountType, serverType string, entryPercent float64, totalData *model.TotalBarData, channel chan bool) {
	defer util.HandlePanic()
	broker := api.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol)
	entryAmount := broker.Cash * entryPercent

	for <-channel {

		pipeline.RefreshPosition(dataModel, broker)
		if dataprocessor.ProcessBarData(dataModel, totalData) {
			log.Printf("curr psar: %f, prev psar: %f, curr bar: %f, prev bar: %f\n",
				dataModel.Signal.CurrentParabolicSar, dataModel.Signal.PreviousParabolicSar, dataModel.Signal.CurrentBar.Close, dataModel.Signal.PreviousBar.Close)
			qty := float64(int(entryAmount / totalData.QuoteData[dataModel.Symbol].AskPrice))
			if signalcatcher.CanEnterLong(dataModel, broker) {
				pipeline.EnterBracketLongPosition(dataModel, totalData, broker, qty)
			} else if signalcatcher.CanEnterShort(dataModel, broker) {
				pipeline.EnterBracketShortPosition(dataModel, totalData, broker, qty)
			}
		}

	}
	log.Printf("%s worker closed", symbol)

}
