package internal

import (
	"github.com/HaoxuanXu/MATradingBot/config"
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
