// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const createChore = `-- name: CreateChore :one
INSERT INTO chores (user_id, title, recurrence_type, recurrence_id, assigned)
VALUES (?, ?, ?, ?, ?)
RETURNING id, user_id, title, recurrence_type, recurrence_id, assigned, created, last_updated
`

type CreateChoreParams struct {
	UserID         int64
	Title          string
	RecurrenceType string
	RecurrenceID   int64
	Assigned       sql.NullString
}

func (q *Queries) CreateChore(ctx context.Context, arg CreateChoreParams) (Chore, error) {
	row := q.db.QueryRowContext(ctx, createChore,
		arg.UserID,
		arg.Title,
		arg.RecurrenceType,
		arg.RecurrenceID,
		arg.Assigned,
	)
	var i Chore
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.RecurrenceType,
		&i.RecurrenceID,
		&i.Assigned,
		&i.Created,
		&i.LastUpdated,
	)
	return i, err
}

const createGroceryItem = `-- name: CreateGroceryItem :one
INSERT INTO groceries
    (user_id, name, completed)
VALUES (?, ?, ?)
RETURNING id, user_id, name, completed, last_updated
`

type CreateGroceryItemParams struct {
	UserID    int64
	Name      string
	Completed bool
}

func (q *Queries) CreateGroceryItem(ctx context.Context, arg CreateGroceryItemParams) (Grocery, error) {
	row := q.db.QueryRowContext(ctx, createGroceryItem, arg.UserID, arg.Name, arg.Completed)
	var i Grocery
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Completed,
		&i.LastUpdated,
	)
	return i, err
}

const createOnceRecurrence = `-- name: CreateOnceRecurrence :one
INSERT INTO chores_recurrence_once (due_date) VALUES (?)
RETURNING id, due_date
`

func (q *Queries) CreateOnceRecurrence(ctx context.Context, dueDate time.Time) (ChoresRecurrenceOnce, error) {
	row := q.db.QueryRowContext(ctx, createOnceRecurrence, dueDate)
	var i ChoresRecurrenceOnce
	err := row.Scan(&i.ID, &i.DueDate)
	return i, err
}

const deleteChore = `-- name: DeleteChore :one
DELETE FROM chores
WHERE id = ? AND user_id = ?
RETURNING id, user_id, title, recurrence_type, recurrence_id, assigned, created, last_updated
`

type DeleteChoreParams struct {
	ID     int64
	UserID int64
}

func (q *Queries) DeleteChore(ctx context.Context, arg DeleteChoreParams) (Chore, error) {
	row := q.db.QueryRowContext(ctx, deleteChore, arg.ID, arg.UserID)
	var i Chore
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.RecurrenceType,
		&i.RecurrenceID,
		&i.Assigned,
		&i.Created,
		&i.LastUpdated,
	)
	return i, err
}

const deleteChoreRecurrenceOnce = `-- name: DeleteChoreRecurrenceOnce :one
DELETE FROM chores_recurrence_once
where id = ?
RETURNING id, due_date
`

func (q *Queries) DeleteChoreRecurrenceOnce(ctx context.Context, id int64) (ChoresRecurrenceOnce, error) {
	row := q.db.QueryRowContext(ctx, deleteChoreRecurrenceOnce, id)
	var i ChoresRecurrenceOnce
	err := row.Scan(&i.ID, &i.DueDate)
	return i, err
}

const deleteGroceryItem = `-- name: DeleteGroceryItem :one
DELETE
FROM groceries
WHERE user_id = ?
  AND id = ?
RETURNING id, user_id, name, completed, last_updated
`

type DeleteGroceryItemParams struct {
	UserID int64
	ID     int64
}

func (q *Queries) DeleteGroceryItem(ctx context.Context, arg DeleteGroceryItemParams) (Grocery, error) {
	row := q.db.QueryRowContext(ctx, deleteGroceryItem, arg.UserID, arg.ID)
	var i Grocery
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Completed,
		&i.LastUpdated,
	)
	return i, err
}

const getMeal = `-- name: GetMeal :one
SELECT user_id, day, breakfast, snack1, lunch, snack2, dinner
FROM meals
WHERE day = ?
  AND user_id = ?
`

type GetMealParams struct {
	Day    time.Time
	UserID int64
}

func (q *Queries) GetMeal(ctx context.Context, arg GetMealParams) (Meal, error) {
	row := q.db.QueryRowContext(ctx, getMeal, arg.Day, arg.UserID)
	var i Meal
	err := row.Scan(
		&i.UserID,
		&i.Day,
		&i.Breakfast,
		&i.Snack1,
		&i.Lunch,
		&i.Snack2,
		&i.Dinner,
	)
	return i, err
}

const listGroceries = `-- name: ListGroceries :many
SELECT id, user_id, name, completed, last_updated
FROM groceries
WHERE user_id = ?
ORDER BY id
`

func (q *Queries) ListGroceries(ctx context.Context, userID int64) ([]Grocery, error) {
	rows, err := q.db.QueryContext(ctx, listGroceries, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Grocery
	for rows.Next() {
		var i Grocery
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Completed,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMeals = `-- name: ListMeals :many
SELECT user_id, day, breakfast, snack1, lunch, snack2, dinner
FROM meals
WHERE user_id = ?
  AND day >= ?
  AND day <= ?
ORDER BY day
`

type ListMealsParams struct {
	UserID   int64
	StartDay time.Time
	EndDay   time.Time
}

func (q *Queries) ListMeals(ctx context.Context, arg ListMealsParams) ([]Meal, error) {
	rows, err := q.db.QueryContext(ctx, listMeals, arg.UserID, arg.StartDay, arg.EndDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Meal
	for rows.Next() {
		var i Meal
		if err := rows.Scan(
			&i.UserID,
			&i.Day,
			&i.Breakfast,
			&i.Snack1,
			&i.Lunch,
			&i.Snack2,
			&i.Dinner,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listOnceChores = `-- name: ListOnceChores :many
SELECT chores.id, chores.user_id, chores.title, chores.recurrence_type, chores.recurrence_id, chores.assigned, chores.created, chores.last_updated, chores_recurrence_once.id, chores_recurrence_once.due_date FROM chores
JOIN chores_recurrence_once  ON chores.recurrence_id = chores_recurrence_once.id
WHERE user_id = ?
`

type ListOnceChoresRow struct {
	Chore                Chore
	ChoresRecurrenceOnce ChoresRecurrenceOnce
}

func (q *Queries) ListOnceChores(ctx context.Context, userID int64) ([]ListOnceChoresRow, error) {
	rows, err := q.db.QueryContext(ctx, listOnceChores, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListOnceChoresRow
	for rows.Next() {
		var i ListOnceChoresRow
		if err := rows.Scan(
			&i.Chore.ID,
			&i.Chore.UserID,
			&i.Chore.Title,
			&i.Chore.RecurrenceType,
			&i.Chore.RecurrenceID,
			&i.Chore.Assigned,
			&i.Chore.Created,
			&i.Chore.LastUpdated,
			&i.ChoresRecurrenceOnce.ID,
			&i.ChoresRecurrenceOnce.DueDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const toggleGroceryItem = `-- name: ToggleGroceryItem :one
UPDATE groceries
SET completed = NOT completed
WHERE user_id = ?
  and id = ?
RETURNING id, user_id, name, completed, last_updated
`

type ToggleGroceryItemParams struct {
	UserID int64
	ID     int64
}

func (q *Queries) ToggleGroceryItem(ctx context.Context, arg ToggleGroceryItemParams) (Grocery, error) {
	row := q.db.QueryRowContext(ctx, toggleGroceryItem, arg.UserID, arg.ID)
	var i Grocery
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Completed,
		&i.LastUpdated,
	)
	return i, err
}

const updateMeals = `-- name: UpdateMeals :one
INSERT OR
REPLACE
INTO meals
    (breakfast, snack1, lunch, snack2, dinner, day, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING user_id, day, breakfast, snack1, lunch, snack2, dinner
`

type UpdateMealsParams struct {
	Breakfast string
	Snack1    string
	Lunch     string
	Snack2    string
	Dinner    string
	Day       time.Time
	UserID    int64
}

func (q *Queries) UpdateMeals(ctx context.Context, arg UpdateMealsParams) (Meal, error) {
	row := q.db.QueryRowContext(ctx, updateMeals,
		arg.Breakfast,
		arg.Snack1,
		arg.Lunch,
		arg.Snack2,
		arg.Dinner,
		arg.Day,
		arg.UserID,
	)
	var i Meal
	err := row.Scan(
		&i.UserID,
		&i.Day,
		&i.Breakfast,
		&i.Snack1,
		&i.Lunch,
		&i.Snack2,
		&i.Dinner,
	)
	return i, err
}
