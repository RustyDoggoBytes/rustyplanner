package main

import (
	"log"
	"net/url"
	"time"
)

var days []string = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

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
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
}

func processWeeklyMealFromForm(form url.Values) []MealPlan {
	var breakfast []string
	var lunch []string
	var snack1 []string
	var snack2 []string
	var dinner []string

	for key, meals := range form {
		if key == "breakfast" {
			breakfast = meals
		} else if key == "snack-1" {
			snack1 = meals
		} else if key == "snack-2" {
			snack2 = meals
		} else if key == "lunch" {
			lunch = meals
		} else if key == "dinner" {
			dinner = meals
		} else {
			log.Printf("ERROR: Unknown key: %s", key)
		}
	}

	newMap := make([]MealPlan, len(days))
	for i, day := range days {
		newMap[i] = MealPlan{
			Day:       day,
			Breakfast: breakfast[i],
			Snack1:    snack1[i],
			Lunch:     lunch[i],
			Snack2:    snack2[i],
			Dinner:    dinner[i],
		}
	}

	return newMap
}
