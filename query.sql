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