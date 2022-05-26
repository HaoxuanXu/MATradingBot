package model

import (
	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

// the ma conditions will track if the 15 minutes 20 day ma and 30 day ma is dropping or rising

func GetDataModel(symbol string, trailsLength int) *DataModel {
	model := DataModel{
		Symbol: symbol,
	}
	model.Trails.ArrayLength = trailsLength
	return &model
}

type DataModel struct {
	Symbol    string
	Condition MAConditions
	CloseData MABarCloseData
	Position  PositionData
	Trails    TrailData
}

type TotalBarData struct {
	Data map[string][]marketdata.Bar
}

type MAConditions struct {
	IsMA20PeriodsDropping           bool
	IsMA20PeriodsPreviouslyDropping bool
	IsMA20PeriodsRising             bool
	IsMA20PeriodsPreviouslyRising   bool
	IsMABelowMA20                   bool
	IsMAAboveMA20                   bool
}

type TrailData struct {
	HWM                 float64
	LongTrailCandidate  float64
	ShortTrailCandidate float64
	LongTrailArray      []float64
	ShortTrailArray     []float64
	ArrayLength         int
	AppliedLongTrail    float64
	AppliedShortTrail   float64
}

type MABarCloseData struct {
	MASupport     float64
	MAResistance  float64
	CurrMAClose   float64
	CurrMA20Close float64
	PrevMA20Close float64
}

type PositionData struct {
	Order            *alpaca.Order
	HasLongPosition  bool
	HasShortPosition bool
	FilledQuantity   float64
	FilledPrice      float64
	CurrentTrail     float64
}

func (position *PositionData) GetPosition(symbol string, broker *internal.AlpacaBroker) {
	positionData, err := broker.GetPosition(symbol)
	if positionData == nil && err != nil {
		position.HasLongPosition = false
		position.HasShortPosition = false
		position.FilledQuantity = 0.0
		position.FilledPrice = 0.0
	}
}
