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

	var longQty float64
	var shortQty float64

	for <-channel {

		pipeline.RefreshPosition(dataModel, broker)
		if dataprocessor.ProcessBarData(dataModel, totalData, isCrypto) {

			if !isCrypto {
				longQty = entryAmount / totalData.StockQuoteData[dataModel.Symbol].AskPrice
				shortQty = float64(int(entryAmount / totalData.StockQuoteData[dataModel.Symbol].AskPrice))
			} else {
				longQty = entryAmount / totalData.CryptoQuoteData[dataModel.Symbol].AskPrice
				shortQty = float64(int(entryAmount / totalData.CryptoQuoteData[dataModel.Symbol].AskPrice))
			}

			if signalcatcher.CanEnterLong(dataModel, broker, isCrypto) && longQty > 0 {
				pipeline.EnterBracketLongPosition(dataModel, totalData, broker, longQty, isCrypto)
			} else if signalcatcher.CanEnterShort(dataModel, broker, isCrypto) && shortQty > 0 {
				pipeline.EnterBracketShortPosition(dataModel, totalData, broker, shortQty, isCrypto)
			}
		}

	}
	log.Printf("%s worker closed", symbol)

}
