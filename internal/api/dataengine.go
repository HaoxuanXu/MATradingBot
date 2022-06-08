package api

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/config"
	"github.com/HaoxuanXu/MATradingBot/util"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

type MarketDataEngine struct {
	client marketdata.Client
}

func GetDataEngine(accountType, serverType string) *MarketDataEngine {

	engine := MarketDataEngine{}
	engine.initialize(accountType, serverType)

	return &engine
}

func (engine *MarketDataEngine) initialize(accountType, serverType string) {
	cred := config.GetCredentials(accountType, serverType)
	engine.client = marketdata.NewClient(
		marketdata.ClientOpts{
			ApiKey:    cred.API_KEY,
			ApiSecret: cred.API_SECRET,
		},
	)
}

func (engine *MarketDataEngine) GetMultiBars(timeframe int, stocks []string) map[string][]marketdata.Bar {

	result, err := engine.client.GetMultiBars(stocks, marketdata.GetBarsParams{
		TimeFrame:  marketdata.NewTimeFrame(timeframe, marketdata.TimeFrameUnit(marketdata.Min)),
		Adjustment: marketdata.Adjustment(marketdata.Raw),
		Start:      util.GetStartTime(time.Now(), 30),
		End:        time.Now(),
		Feed:       "sip",
	})
	if err != nil {
		log.Println(err)
	}

	return result
}

func (engine *MarketDataEngine) GetMultiCryptoBars(timeframe int, cryptos []string) map[string][]marketdata.CryptoBar {
	result, err := engine.client.GetCryptoMultiBars(cryptos, marketdata.GetCryptoBarsParams{
		TimeFrame: marketdata.NewTimeFrame(timeframe, marketdata.TimeFrameUnit(marketdata.Min)),
		Start:     util.GetStartTime(time.Now(), 30),
		End:       time.Now(),
	})
	if err != nil {
		log.Println(err)
	}

	return result
}

func (engine *MarketDataEngine) GetLatestMultiTrades(assets []string) map[string]marketdata.Trade {
	result, err := engine.client.GetLatestTrades(assets)
	if err != nil {
		log.Println(err)
	}
	return result
}

func (engine *MarketDataEngine) GetLatestMultiQuotes(assets []string) map[string]marketdata.Quote {

	result, err := engine.client.GetLatestQuotes(assets)
	if err != nil {
		log.Println(err)
	}

	return result
}

func (engine *MarketDataEngine) GetLatestMultiCryptoQuotes(cryptos []string) map[string]marketdata.CryptoQuote {
	result, err := engine.client.GetLatestCryptoQuotes(cryptos, "FTXU")
	if err != nil {
		log.Println(err)
	}

	return result
}
