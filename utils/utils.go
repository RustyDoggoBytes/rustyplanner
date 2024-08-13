package utils

import (
	"os"
	"time"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}

func GetMondayOfCurrentWeek(now time.Time) time.Time {
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

func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func FormatMonthDay(date time.Time) string {
	return date.Format("01/02")
}
