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
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema"
	"github.com/HaoxuanXu/MATradingBot/strats/macdpsar200ema/model"
	"github.com/HaoxuanXu/MATradingBot/util"
)

func main() {
	// read from flag
	yamlFileName := flag.String("config", "production-paper-account.yml", "this yml config file for the application")
	flag.Parse()

	yamlConfig := util.ReadYAMLFile(db.MapYAMLConfigPath(*yamlFileName))
	accountType := fmt.Sprintf("%s", yamlConfig["accounttype"])
	serverType := fmt.Sprintf("%s", yamlConfig["servertype"])
	entryPercent, _ := strconv.ParseFloat(fmt.Sprintf("%v", yamlConfig["entrypercent"]), 64)

	var fastTotalData model.TotalBarData
	var slowTotalData model.TotalBarData
	stocks := config.Stocks

	workerEntryPercent := entryPercent / float64(len(stocks))

	// set up logging
	logFile := logging.SetLogging()
	log.Printf("Number of Workers: %d\n", len(stocks))

	dataEngine := api.GetDataEngine(accountType, serverType)
	broker := api.GetBroker(accountType, serverType)

	clock, _ := broker.GetClock()
	if !clock.IsOpen {
		log.Printf("Market closed currently. Wait for %.2f minutes\n", time.Until(clock.NextOpen).Minutes())
		time.Sleep(time.Until(clock.NextOpen))
	}

	// create channel map
	stockChanMap := channel.CreateMap(stocks)

	// start workers
	for _, stock := range stocks {
		log.Printf("Starting worker for %s trading\n", stock)
		go macdpsar200ema.MACDPSar200EMAStrategy(stock, accountType, serverType, workerEntryPercent, &fastTotalData, &slowTotalData, stockChanMap.Map[stock])
	}

	// start main loop
	clock, _ = broker.GetClock()
	log.Println("Start main loop...")
	for time.Until(clock.NextClose) > 0 {

		fastBarData := dataEngine.GetMultiBars(1, stocks)
		if len(fastBarData) > 0 {
			fastTotalData.StockBarData = fastBarData
		}
		slowBarData := dataEngine.GetMultiBars(5, stocks)
		if len(slowBarData) > 0 {
			slowTotalData.StockBarData = slowBarData
		}

		quoteData := dataEngine.GetLatestMultiQuotes(stocks)
		slowTotalData.StockQuoteData = quoteData
		fastTotalData.StockQuoteData = quoteData
		stockChanMap.TriggerWorkers()
		time.Sleep(time.Second)
	}

	stockChanMap.CloseWorkers()
	time.Sleep(5 * time.Second) // wait for all worker routines to close
	stockChanMap.CloseChannels()
	logFile.Close()

}
