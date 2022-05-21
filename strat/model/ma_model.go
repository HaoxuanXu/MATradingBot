package model

import (
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
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
	CurrMA20Bar alpaca.Bar
	CurrMA30Bar alpaca.Bar
	PrevMA20Bar alpaca.Bar
	PrevMA30Bar alpaca.Bar
}
