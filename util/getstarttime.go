package util

import "time"

func GetStartTime(endTime time.Time, days int) time.Time {
	startTime := endTime

	for days > 0 {
		startTime = startTime.AddDate(0, 0, -1)
		if startTime.Weekday() != time.Saturday || startTime.Weekday() != time.Sunday {
			days--
		}
	}

	return startTime
}
