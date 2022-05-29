package main

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/config"
	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/HaoxuanXu/MATradingBot/strat"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
)

func main() {

	var accountType string
	var serverType string
	var entryPercent float64

	var totalData model.TotalBarData
	assets := config.Assets
	channelMapper := make(map[string](chan bool))
	dataEngine := internal.GetDataEngine(accountType, serverType)
	broker := internal.GetBroker(accountType, serverType)

	for _, asset := range assets {
		channelMapper[asset] = make(chan bool)
	}
	for _, asset := range assets {
		log.Printf("Starting worker for %s trading\n", asset)
		go strat.MATradingStrategy(asset, accountType, serverType, entryPercent, &totalData, channelMapper[asset])
	}

	if !broker.Clock.IsOpen {
		log.Printf("Wait for %.2f minutes till the market opens\n", time.Until(broker.Clock.NextOpen).Minutes())
		time.Sleep(time.Until(broker.Clock.NextOpen))
	}

	for broker.Clock.IsOpen {
		totalData.Data = dataEngine.GetMultiBars(1, assets)
		time.Sleep(time.Minute)
	}
}
