package db

import (
	"context"
	"database/sql"
	sqlc "rustydoggobytes/planner/sqlc_generated"
	"time"
)

type MealPlan struct {
	Date      time.Time
	Breakfast string
	Snack1    string
	Lunch     string
	Snack2    string
	Dinner    string
}

type Repository struct {
	db      *sql.DB
	ctx     context.Context
	queries *sqlc.Queries
}

func NewRepository(ctx context.Context, db *sql.DB, schema string) (*Repository, error) {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, err
	}

	queries := sqlc.New(db)
	return &Repository{ctx: ctx, db: db, queries: queries}, nil
}

func (r *Repository) GetMealPlanByDate(userID int64, startDate, endDate time.Time) ([]MealPlan, error) {
	startDate = startDate.Truncate(24 * time.Hour)
	endDate = endDate.Truncate(24 * time.Hour)

	params := sqlc.ListMealsParams{
		UserID:   userID,
		StartDay: startDate,
		EndDay:   endDate,
	}
	dbMeals, err := r.queries.ListMeals(r.ctx, params)
	if err != nil {
		return nil, err
	}

	var meals []MealPlan
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		meal := MealPlan{Date: date}
		for _, dbMeal := range dbMeals {
			if dbMeal.Day == date {
				meal = MealPlan{
					Date:      dbMeal.Day,
					Breakfast: dbMeal.Breakfast,
					Snack1:    dbMeal.Snack1,
					Snack2:    dbMeal.Snack2,
					Lunch:     dbMeal.Lunch,
					Dinner:    dbMeal.Dinner,
				}
				break
			}
		}
		meals = append(meals, meal)
	}

	return meals, nil
}

func (r *Repository) UpdateMealPlan(userID int64, meal MealPlan) error {
	date := meal.Date.Truncate(24 * time.Hour)

	params := sqlc.UpdateMealsParams{
		UserID:    userID,
		Day:       date,
		Breakfast: meal.Breakfast,
		Snack1:    meal.Snack1,
		Snack2:    meal.Snack2,
		Lunch:     meal.Lunch,
		Dinner:    meal.Dinner,
	}
	_, err := r.queries.UpdateMeals(r.ctx, params)
	if err != nil {
		return err
	}
	return nil
}
