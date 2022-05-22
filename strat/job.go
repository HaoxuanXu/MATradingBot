package strat

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/HaoxuanXu/MATradingBot/strat/dataprocessor"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
)

func MATradingStrategy(symbol, accountType, serverType string, totalData *model.TotalBarData) {
	broker := internal.GetBroker(accountType, serverType)
	dataModel := model.GetDataModel(symbol, 30)

	for !broker.Clock.IsOpen {
		log.Printf("Wait for %.2f minutes till the market opens\n", time.Until(broker.Clock.NextOpen).Minutes())
	}

	for broker.Clock.IsOpen {
		dataprocessor.ProcessBarData(&dataModel, totalData)

	}
}
