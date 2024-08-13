package components

import (
	"rustydoggobytes/planner/db"
	"time"
)

type PageData struct {
	WeekStart    time.Time
	WeekEnd      time.Time
	PreviousWeek time.Time
	NextWeek     time.Time
	Meals        []db.MealPlan
	FormData     map[string][]string
}
