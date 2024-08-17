-- name: ListMeals :many
SELECT *
FROM meals
WHERE user_id = ?
  AND day >= @start_day
  AND day <= @end_day
ORDER BY day;

-- name: GetMeal :one
SELECT *
FROM meals
WHERE day = ?
  AND user_id = ?;

-- name: UpdateMeals :one
INSERT OR
REPLACE
INTO meals
    (breakfast, snack1, lunch, snack2, dinner, day, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListGroceries :many
SELECT *
FROM groceries
WHERE user_id = ?
ORDER BY id;

-- name: DeleteGroceryItem :one
DELETE
FROM groceries
WHERE user_id = ?
  AND id = ?
RETURNING *;

-- name: ToggleGroceryItem :one
UPDATE groceries
SET completed = NOT completed
WHERE user_id = ?
  and id = ?
RETURNING *;

-- name: CreateGroceryItem :one
INSERT INTO groceries
    (user_id, name, completed)
VALUES (?, ?, ?)
RETURNING *;

-- name: ListOnceChores :many
SELECT sqlc.embed(chores), sqlc.embed(chores_recurrence_once) FROM chores
JOIN chores_recurrence_once  ON chores.recurrence_id = chores_recurrence_once.id
WHERE user_id = ?;

-- name: CreateChore :one
INSERT INTO chores (user_id, title, recurrence_type, recurrence_id, assigned)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: CreateOnceRecurrence :one
INSERT INTO chores_recurrence_once (due_date) VALUES (?)
RETURNING *;

-- name: DeleteChore :one
DELETE FROM chores
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteChoreRecurrenceOnce :one
DELETE FROM chores_recurrence_once
where id = ?
RETURNING *;
