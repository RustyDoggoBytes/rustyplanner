package db

import (
	sqlc "rustydoggobytes/planner/sqlc_generated"
	"sort"
	"strconv"
	"time"
)

type GroceryItem struct {
	ID          int64
	Name        string
	Completed   bool
	LastUpdated time.Time
}

func (r *Repository) ListGroceryItems(userID int64) ([]GroceryItem, error) {
	dbGroceries, err := r.queries.ListGroceries(r.ctx, userID)
	if err != nil {
		return nil, err
	}

	groceries := make([]GroceryItem, len(dbGroceries))
	for i, grocery := range dbGroceries {
		groceries[i] = *mapToGroceryItem(grocery)
	}

	sort.Slice(groceries, func(i, j int) bool {
		if groceries[i].Completed != groceries[j].Completed {
			return !groceries[i].Completed
		}

		return groceries[i].LastUpdated.After(groceries[j].LastUpdated)
	})

	return groceries, nil
}

func (r *Repository) CreateGroceryItem(userID int64, name string) (*GroceryItem, error) {
	params := sqlc.CreateGroceryItemParams{
		UserID:    userID,
		Name:      name,
		Completed: false,
	}
	groceryItem, err := r.queries.CreateGroceryItem(r.ctx, params)
	if err != nil {
		return nil, err
	}

	return mapToGroceryItem(groceryItem), nil
}

func (r *Repository) ToggleGroceryItem(userID int64, id string) (*GroceryItem, error) {
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	params := sqlc.ToggleGroceryItemParams{UserID: userID, ID: id_int}
	item, err := r.queries.ToggleGroceryItem(r.ctx, params)
	if err != nil {
		return nil, err
	}

	return mapToGroceryItem(item), nil
}

func (r *Repository) DeleteGroceryItem(userID int64, id string) error {
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	params := sqlc.DeleteGroceryItemParams{UserID: userID, ID: id_int}

	_, err = r.queries.DeleteGroceryItem(r.ctx, params)
	return err
}

func mapToGroceryItem(grocery sqlc.Grocery) *GroceryItem {
	return &GroceryItem{
		ID:          grocery.ID,
		Name:        grocery.Name,
		Completed:   grocery.Completed,
		LastUpdated: grocery.LastUpdated,
	}
}
