package internal

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/config"
	"github.com/HaoxuanXu/MATradingBot/util"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

type MarketDataEngine struct {
	client marketdata.Client
	Assets []string
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
	engine.Assets = config.Assets
}

func (engine *MarketDataEngine) GetMultiBars(timeframe, days int) map[string][]marketdata.Bar {

	result, err := engine.client.GetMultiBars(engine.Assets, marketdata.GetBarsParams{
		TimeFrame:  marketdata.NewTimeFrame(timeframe, marketdata.TimeFrameUnit(marketdata.Min)),
		Adjustment: marketdata.Adjustment(marketdata.Raw),
		Start:      util.GetStartTime(time.Now(), days),
		End:        time.Now(),
	})
	if err != nil {
		log.Println(err)
	}

	return result
}
