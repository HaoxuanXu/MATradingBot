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

	var totalData model.TotalBarData
	stocks := config.Stocks
	cryptos := config.Crypto

	workerEntryPercent := entryPercent / float64(len(stocks)+len(cryptos))

	// set up logging
	_ = logging.SetLogging()

	dataEngine := api.GetDataEngine(accountType, serverType)
	broker := api.GetBroker(accountType, serverType)

	// create channel map
	stockChanMap := channel.CreateMap(stocks)
	cryptoChanMap := channel.CreateMap(cryptos)

	// start workers
	for _, stock := range stocks {
		log.Printf("Starting worker for %s trading\n", stock)
		go macdpsar200ema.MACDPSar200EMAStrategy(stock, accountType, serverType, workerEntryPercent, &totalData, false, stockChanMap.Map[stock])
	}

	for _, crypto := range cryptos {
		log.Printf("Starting worker for %s trading\n", crypto)
		go macdpsar200ema.MACDPSar200EMAStrategy(crypto, accountType, serverType, workerEntryPercent, &totalData, true, cryptoChanMap.Map[crypto])
	}

	// start main loop
	log.Println("Start main loop...")
	for {
		clock, _ := broker.GetClock()
		if clock.IsOpen {
			barData := dataEngine.GetMultiBars(15, stocks)
			if len(barData) > 0 {
				totalData.StockBarData = barData
			}
			totalData.StockQuoteData = dataEngine.GetLatestMultiQuotes(stocks)
			stockChanMap.TriggerWorkers()
		}
		cryptoBarData := dataEngine.GetMultiCryptoBars(15, cryptos)
		if len(cryptoBarData) > 0 {
			totalData.CryptoBarData = cryptoBarData
		}
		totalData.CryptoQuoteData = dataEngine.GetLatestMultiCryptoQuotes(cryptos)
		cryptoChanMap.TriggerWorkers()
		time.Sleep(time.Minute)
	}
}
