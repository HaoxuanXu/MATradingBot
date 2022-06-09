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

func MACDPSar200EMAStrategy(symbol, accountType, serverType string, entryPercent float64, totalData *model.TotalBarData, isCrypto bool, channel chan bool) {
	defer util.HandlePanic(symbol)
	broker := api.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol)
	entryAmount := broker.Cash * entryPercent

	var qty float64

	for <-channel {

		pipeline.RefreshPosition(dataModel, broker)
		if dataprocessor.ProcessBarData(dataModel, totalData, isCrypto) {

			if !isCrypto {
				qty = float64(int(entryAmount / totalData.StockQuoteData[dataModel.Symbol].AskPrice))
			} else {
				qty = float64(int(entryAmount / totalData.CryptoQuoteData[dataModel.Symbol].AskPrice))
			}

			if signalcatcher.CanEnterLong(dataModel, broker, isCrypto) && qty > 0 {
				pipeline.EnterBracketLongPosition(dataModel, totalData, broker, qty, isCrypto)
			} else if signalcatcher.CanEnterShort(dataModel, broker, isCrypto) && qty > 0 {
				pipeline.EnterBracketShortPosition(dataModel, totalData, broker, qty, isCrypto)
			}
		}

	}
	log.Printf("%s worker closed", symbol)

}
