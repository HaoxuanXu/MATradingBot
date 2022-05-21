package model

import "github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"

// the ma conditions will track if the 15 minutes 20 day ma and 30 day ma is dropping or rising

type TotalBarData struct {
	Data map[string][]marketdata.Bar
}

type MAConditions struct {
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
}
