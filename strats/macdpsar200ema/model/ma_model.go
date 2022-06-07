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
	BarData   map[string][]marketdata.Bar
	QuoteData map[string]marketdata.Quote
}

type SignalData struct {
	CurrentEMA200Period  float64
	CurrentParabolicSar  float64
	CurrentBar           marketdata.Bar
	CurrentClose         float64
	PreviousClose        float64
	PreviousParabolicSar float64
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
