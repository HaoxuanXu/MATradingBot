package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/HaoxuanXu/MATradingBot/config"
	"github.com/HaoxuanXu/MATradingBot/db"
	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/internal/channel"
	"github.com/HaoxuanXu/MATradingBot/internal/logging"
	"github.com/HaoxuanXu/MATradingBot/strat"
	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/HaoxuanXu/MATradingBot/strat/tools"
	"github.com/HaoxuanXu/MATradingBot/util"
)

func main() {
	// read from flag
	yamlFileName := flag.String("config", "production-paper-account.yml", "this yml config file for the application")
	flag.Parse()

	yamlConfig := util.ReadYAMLFile(db.MapYAMLConfigPath(*yamlFileName))

	var totalData model.TotalBarData
	assets := config.Assets

	accountType := fmt.Sprintf("%s", yamlConfig["accounttype"])
	serverType := fmt.Sprintf("%s", yamlConfig["servertype"])
	totalEntryPercent, _ := strconv.ParseFloat(fmt.Sprintf("%v", yamlConfig["entrypercent"]), 64)
	workerEntryPercent := totalEntryPercent / float64(len(assets))

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
		go strat.MATradingStrategy(asset, accountType, serverType, workerEntryPercent, &totalData, chanMap.Map[asset])
	}

	// start main loop
	log.Println("Start main loop...")
	broker.Clock, _ = broker.Client.GetClock()
	for time.Until(broker.Clock.NextClose) > 0 {
		barData := dataEngine.GetMultiBars(1, assets)
		if len(barData) > 0 {
			totalData.BarData = barData
		}
		totalData.TradeData = dataEngine.GetLatestMultiTrades(assets)
		for key := range totalData.BarData {
			totalData.BarData[key] = tools.Reverse(totalData.BarData[key])
		}
		chanMap.TriggerWorkers()
		time.Sleep(time.Minute)
	}
	// close operation when the market is closed
	log.Println("Shutting down workers...")
	chanMap.CloseWorkers()
	log.Println("Closing channels...")
	chanMap.CloseChannels()
	logFile.Close()
}
