package macdpsar200ema

import (
	"log"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/dataprocessor"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/pipeline"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/signalcatcher"
)

func MACDPSar200EMAStrategy(symbol, accountType, serverType string, entryPercent float64, totalData *model.TotalBarData, channel chan bool) {
	// defer util.HandlePanic(symbol)
	broker := api.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol)
	entryAmount := broker.Cash * entryPercent

	var qty float64

	for <-channel {

		pipeline.RefreshPosition(dataModel, broker)
		if dataprocessor.ProcessBarData(dataModel, totalData) {
			qty = float64(int(entryAmount / dataModel.Signal.Quote.AskPrice))

			if signalcatcher.CanEnterLong(dataModel, broker) && qty > 0 {
				pipeline.EnterTrailingStopLongPosition(dataModel, broker, qty)
			} else if signalcatcher.CanEnterShort(dataModel, broker) && qty > 0 {
				pipeline.EnterTrailingStopShortPosition(dataModel, broker, qty)
			}
		}

	}
	log.Printf("%s worker closed", symbol)

}
