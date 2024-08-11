package main

import (
	"time"
)

func getMondayOfCurrentWeek(now time.Time) time.Time {
	weekday := now.Weekday()
	var monday time.Time
	if weekday == time.Sunday {
		// If today is Sunday, we need to go back 6 days
		monday = now.AddDate(0, 0, -6)
	} else {
		// Otherwise, we go back (weekday - 1) days
		monday = now.AddDate(0, 0, -int(weekday)+1)
	}
	return monday.Truncate(24 * time.Hour)
}
