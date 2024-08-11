package db

import "time"

type GroceryItem struct {
	ID          int64
	Name        string
	Completed   bool
	LastUpdated time.Time
}
