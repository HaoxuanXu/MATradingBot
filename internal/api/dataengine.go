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
	engine := &MarketDataEngine{}
	engine.initialize(accountType, serverType)
	return engine
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

func (engine *MarketDataEngine) GetMultiBars(timeframe int, assets []string) map[string][]marketdata.Bar {

	result, err := engine.client.GetMultiBars(assets, marketdata.GetBarsParams{
		TimeFrame:  marketdata.NewTimeFrame(timeframe, marketdata.TimeFrameUnit(marketdata.Min)),
		Adjustment: marketdata.Adjustment(marketdata.Raw),
		Start:      util.GetStartTime(time.Now(), 1),
		End:        time.Now(),
		Feed:       "sip",
	})
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
