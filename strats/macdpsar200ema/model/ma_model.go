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
	Shortable           bool
	CurrentBarTimestamp time.Time
	Signal              SignalData
	Position            PositionData
}

type TotalBarData struct {
	StockBarData   map[string][]marketdata.Bar
	StockQuoteData map[string]marketdata.Quote
}

type SignalData struct {
	EMA200Periods         []float64
	ParabolicSars         []float64
	StochK                []float64
	StochD                []float64
	StochOversold         bool
	StochOverbought       bool
	RSI                   []float64
	BXTrenderLongTerm     []float64
	BXTrenderShortTerm    []float64
	Bars                  []marketdata.Bar
	Quote                 marketdata.Quote
	TrailingStopLossLong  float64
	TrailingStopLossShort float64
	Macds                 []float64
	MacdSignals           []float64
	FibonacciLow          []float64
	FibonacciHigh         []float64
	SwingLow              []float64
	SwingHigh             []float64
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
