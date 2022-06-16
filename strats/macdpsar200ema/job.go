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

func MACDPSar200EMAStrategy(
	symbol, accountType, serverType string,
	entryPercent float64,
	fastTotalData, slowTotalData *model.TotalBarData,
	channel chan bool,
) {
	defer util.HandlePanic(symbol)
	broker := api.GetBroker(accountType, serverType)
	fastDataModel := model.GetDataModel(symbol)
	slowDataModel := model.GetDataModel(symbol)
	entryAmount := broker.Cash * entryPercent

	var prevTimeFrameDirectionAlignedLong bool
	var prevTimeFrameDirectionAlignedShort bool

	var qty float64

	for <-channel {

		pipeline.RefreshPosition(fastDataModel, broker)
		if dataprocessor.ProcessBarData(fastDataModel, fastTotalData) {
			dataprocessor.ProcessBarData(slowDataModel, slowTotalData)
			qty = float64(int(entryAmount / fastDataModel.Signal.Quote.AskPrice))

			if signalcatcher.CanEnterLong(fastDataModel) &&
				signalcatcher.CanEnterLong(slowDataModel) &&
				!broker.HasLongPosition && !broker.HasShortPosition &&
				qty > 0 && !prevTimeFrameDirectionAlignedLong {
				pipeline.EnterTrailingStopLongPosition(fastDataModel, broker, qty)
			} else if signalcatcher.CanEnterShort(fastDataModel) &&
				signalcatcher.CanEnterShort(slowDataModel) &&
				!broker.HasLongPosition && !broker.HasShortPosition && qty > 0 &&
				!prevTimeFrameDirectionAlignedShort {
				pipeline.EnterTrailingStopShortPosition(fastDataModel, broker, qty)
			}
		}

		// record slow time frame fast time frame alignment
		if signalcatcher.CanEnterLong(slowDataModel) && signalcatcher.CanEnterLong(fastDataModel) {
			prevTimeFrameDirectionAlignedLong = true
		} else {
			prevTimeFrameDirectionAlignedLong = false
		}

		if signalcatcher.CanEnterShort(fastDataModel) && signalcatcher.CanEnterShort(slowDataModel) {
			prevTimeFrameDirectionAlignedShort = true
		} else {
			prevTimeFrameDirectionAlignedShort = false
		}

	}
	log.Printf("%s worker closed", symbol)

}
