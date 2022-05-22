package model

import (
	"github.com/HaoxuanXu/MATradingBot/internal"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

// the ma conditions will track if the 15 minutes 20 day ma and 30 day ma is dropping or rising

func GetDataModel(symbol string, trailsLength int) DataModel {
	model := DataModel{
		Symbol: symbol,
	}
	model.Trails.DataLength = trailsLength
	return model
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
	HWM                float64
	Trail              float64
	IsMA20DaysDropping bool
	IsMA30DaysDropping bool
	IsMA20DaysRising   bool
	IsMA30DaysRising   bool
	IsMA20AboveMA30    bool
	IsMA20BelowMA30    bool
}

type MABarCloseData struct {
	CurrMA20Close float64
	CurrMA30Close float64
	PrevMA20Close float64
	PrevMA30Close float64
}

type PositionData struct {
	HasLongPosition  bool
	HasShortPosition bool
	FilledQuantity   float64
	FilledPrice      float64
	CurrentTrail     float64
}

type TrailData struct {
	TrailDataArray []float64
	DataLength     int
}

func (position *PositionData) GetPosition(symbol string, broker *internal.AlpacaBroker) {
	positionData := broker.GetPosition(symbol)
	if positionData == nil {
		position.HasLongPosition = false
		position.HasShortPosition = false
		position.FilledQuantity = 0.0
		position.FilledPrice = 0.0
	}
}
