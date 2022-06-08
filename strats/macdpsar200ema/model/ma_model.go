package model

import (
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

// the ma conditions will track if the 15 minutes 20 day ma and 30 day ma is dropping or rising

func GetDataModel(symbol string) *DataModel {
	model := DataModel{
		Symbol: symbol,
	}
	return &model
}

type DataModel struct {
	Symbol              string
	CurrentBarTimestamp time.Time
	Signal              SignalData
	Position            PositionData
}

type TotalBarData struct {
	StockBarData    map[string][]marketdata.Bar
	StockQuoteData  map[string]marketdata.Quote
	CryptoBarData   map[string][]marketdata.CryptoBar
	CryptoQuoteData map[string]marketdata.CryptoQuote
}

type SignalData struct {
	CurrentEMA200Period  float64
	CurrentParabolicSar  float64
	PreviousParabolicSar float64
	CurrentBar           marketdata.Bar
	PreviousBar          marketdata.Bar
	CurrentCryptoBar     marketdata.CryptoBar
	PreviousCryptoBar    marketdata.CryptoBar
	CurrentMacd          float64
	CurrentMacdSignal    float64
}

type PositionData struct {
	Order            alpaca.Order
	HasOrder         bool
	HasLongPosition  bool
	HasShortPosition bool
	FilledQuantity   float64
	FilledPrice      float64
}

func (position *PositionData) GetPosition(symbol string, broker *api.AlpacaBroker) {
	positionData, err := broker.GetPosition(symbol)
	if positionData == nil && err != nil {
		position.HasLongPosition = false
		position.HasShortPosition = false
		position.FilledQuantity = 0.0
		position.FilledPrice = 0.0
	}
}
