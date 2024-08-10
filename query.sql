-- name: ListMeals :many
SELECT *
FROM meals
WHERE user_id = ?
  AND day >= @start_day
  AND day <= @end_day
ORDER BY day;

-- name: GetMeal :one
SELECT *
from meals
where day = ?
  AND user_id = ?;

-- name: UpdateMeals :one
INSERT OR
REPLACE INTO meals
    (breakfast, snack1, lunch, snack2, dinner, day, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;