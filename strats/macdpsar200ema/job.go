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
	defer util.HandlePanic(symbol)
	broker := api.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol)
	entryAmount := broker.Cash * entryPercent

	for <-channel {

		pipeline.RefreshPosition(dataModel, broker)
		if dataprocessor.ProcessBarData(dataModel, totalData) {
			// log.Printf("%s: curr: %.3f, prev: %.3f\n", dataModel.Symbol, dataModel.Signal.RecentParabolicSarDiff[len(dataModel.Signal.RecentParabolicSarDiff)-1],
			// 	dataModel.Signal.RecentParabolicSarDiff[len(dataModel.Signal.RecentParabolicSarDiff)-2])
			longQty := entryAmount / totalData.QuoteData[dataModel.Symbol].AskPrice
			shortQty := float64(int(entryAmount / totalData.QuoteData[dataModel.Symbol].AskPrice))
			if signalcatcher.CanEnterLong(dataModel, broker) && longQty > 0 {
				pipeline.EnterBracketLongPosition(dataModel, totalData, broker, longQty)
			} else if signalcatcher.CanEnterShort(dataModel, broker) && shortQty > 0 {
				pipeline.EnterBracketShortPosition(dataModel, totalData, broker, shortQty)
			}
		}

	}
	log.Printf("%s worker closed", symbol)

}
