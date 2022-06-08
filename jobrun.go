package main

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/config"
	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/internal/channel"
	"github.com/HaoxuanXu/MATradingBot/internal/logging"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
)

func Run(accountType, serverType string, entryPercent float64) {
	// read from flag

	var totalData model.TotalBarData
	assets := config.Assets

	workerEntryPercent := entryPercent / float64(len(assets))

	// set up logging
	logFile := logging.SetLogging()

	dataEngine := api.GetDataEngine(accountType, serverType)
	broker := api.GetBroker(accountType, serverType)

	// create channel map
	chanMap := channel.CreateMap(assets)

	if !broker.Clock.IsOpen {
		log.Printf("Wait for %.2f minutes till the market opens\n", time.Until(broker.Clock.NextOpen).Minutes())
		time.Sleep(time.Until(broker.Clock.NextOpen))
	}

	// start workers
	for _, asset := range assets {
		log.Printf("Starting worker for %s trading\n", asset)
		go macdpsar200ema.MACDPSar200EMAStrategy(asset, accountType, serverType, workerEntryPercent, &totalData, chanMap.Map[asset])
	}

	// start main loop
	log.Println("Start main loop...")
	broker.Clock, _ = broker.Client.GetClock()
	for time.Until(broker.Clock.NextClose) > 0 {
		barData := dataEngine.GetMultiBars(30, assets)
		if len(barData) > 0 {
			totalData.BarData = barData
		}
		totalData.QuoteData = dataEngine.GetLatestMultiQuotes(assets)
		chanMap.TriggerWorkers()
		time.Sleep(time.Minute)
	}
	// close operation when the market is closed
	log.Println("Shutting down workers...")
	chanMap.CloseWorkers()
	time.Sleep(5 * time.Second)
	log.Println("Closing channels...")
	chanMap.CloseChannels()
	logFile.Close()
}
