package util

import (
	"time"
)

func GetStartTime(endTime time.Time, numDays int) time.Time {
	days := 0
	startTime := endTime

	for days < numDays {
		startTime = startTime.AddDate(0, 0, -1)
		if startTime.Weekday() != time.Saturday || startTime.Weekday() != time.Sunday {
			days++
		}
	}
	return startTime
}
