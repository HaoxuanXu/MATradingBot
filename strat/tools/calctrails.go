package tools

import (
	"math"

	"github.com/HaoxuanXu/MATradingBot/strat/model"
	"github.com/montanaflynn/stats"
)

func CalcAppliedTrails(model *model.DataModel, data *model.TotalBarData) {

	// reverse the bar data so it starts from the oldest time
	var longHWM float64
	var shortHWM float64
	var longTrailCandidate float64
	var shortTrailCandidate float64
	var currMA20Close float64
	var longTrailsArray []float64
	var shortTrailsArray []float64
	barData := Reverse(data.BarData[model.Symbol])

	for startIndex := 19; startIndex < len(barData); startIndex++ {
		if longHWM == 0 {
			longHWM = barData[startIndex].High
		}
		if shortHWM == 0 {
			shortHWM = barData[startIndex].Low
		}
		var subArray []float64
		for _, val := range barData[startIndex-19 : startIndex] {
			subArray = append(subArray, val.Close)
		}
		currMA20Close, _ = stats.Mean(subArray[startIndex-19 : startIndex])
		currentBar := barData[startIndex]

		if currentBar.Low > currMA20Close {
			if currentBar.High < longHWM {
				longTrailCandidate = math.Max(longTrailCandidate, longHWM-currentBar.Low)
			} else if currentBar.High > longHWM {
				if longTrailCandidate > 0 {
					longTrailsArray = append(longTrailsArray, longTrailCandidate)
				}
				longTrailCandidate = 0.0
				longHWM = currentBar.High
			}
		} else if currentBar.High < currMA20Close {
			if currentBar.Low > shortHWM {
				shortTrailCandidate = math.Max(shortTrailCandidate, currentBar.High-shortHWM)
			} else if currentBar.Low < shortHWM {
				if shortTrailCandidate > 0.0 {
					shortTrailsArray = append(shortTrailsArray, shortTrailCandidate)
				}
				shortTrailCandidate = 0.0
				shortHWM = currentBar.Low
			}
		}
	}

	shortTrailApplied, _ := stats.Mean(shortTrailsArray)
	longTrailApplied, _ := stats.Mean(longTrailsArray)

	model.Trails.AppliedLongTrail = longTrailApplied
	model.Trails.AppliedShortTrail = shortTrailApplied
}
