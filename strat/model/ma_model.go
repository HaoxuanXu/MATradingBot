package model

import (
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

// the ma conditions will track if the 15 minutes 20 day ma and 30 day ma is dropping or rising
type MAConditions struct {
	HasLongPosition    bool
	HasShortPosition   bool
	IsMA20DaysDropping bool
	IsMA30DaysDropping bool
	IsMA20DaysRising   bool
	IsMA30DaysRising   bool
	IsMA20AboveMA30    bool
	IsMA20BelowMA30    bool
}

type MABarData struct {
	CurrMA20Bar marketdata.Bar
	CurrMA30Bar marketdata.Bar
	PrevMA20Bar marketdata.Bar
	PrevMA30Bar marketdata.Bar
}
