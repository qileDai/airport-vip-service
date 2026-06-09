package services

import (
	"time"
)

func parseDateRange(startDate, endDate string) (time.Time, time.Time) {
	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			start = time.Now().AddDate(0, -1, 0)
		}
	} else {
		start = time.Now().AddDate(0, -1, 0)
	}

	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			end = time.Now()
		}
	} else {
		end = time.Now()
	}

	return start, end
}
