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

	// var prevLongDecisionSlow bool
	// var prevLongDecisionFast bool
	// var prevShortDecisionSlow bool
	// var prevShortDecisionFast bool

	var qty float64

	for <-channel {
		pipeline.RefreshPosition(fastDataModel, broker)
		if dataprocessor.ProcessBarData(fastDataModel, fastTotalData) {
			dataprocessor.ProcessBarData(slowDataModel, slowTotalData)
			qty = float64(int(entryAmount / fastDataModel.Signal.Quote.AskPrice))

			if signalcatcher.CanEnterLongFastPeriod(fastDataModel) &&
				signalcatcher.CanEnterLongSlowPeriod(slowDataModel) &&
				!broker.HasLongPosition && !broker.HasShortPosition &&
				qty > 0 {
				pipeline.EnterBracketLongPosition(fastDataModel, broker, qty)
			} else if signalcatcher.CanEnterShortFastPeriod(fastDataModel) &&
				signalcatcher.CanEnterShortSlowPeriod(slowDataModel) &&
				!broker.HasLongPosition && !broker.HasShortPosition && qty > 0 {
				pipeline.EnterBracketShortPosition(fastDataModel, broker, qty)
			}
		}

		// record slow time frame fast time frame alignment
		// prevLongDecisionFast = signalcatcher.CanEnterLongFastPeriod(fastDataModel)
		// prevLongDecisionSlow = signalcatcher.CanEnterLongSlowPeriod(slowDataModel)
		// prevShortDecisionFast = signalcatcher.CanEnterShortFastPeriod(fastDataModel)
		// prevShortDecisionSlow = signalcatcher.CanEnterShortSlowPeriod(slowDataModel)

	}
	log.Printf("%s worker closed", symbol)

}
