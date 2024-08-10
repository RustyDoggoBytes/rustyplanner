package main

import (
	"context"
	"database/sql"
	"log"
	sqlc "rustydoggobytes/planner/sqlc_generated"
	"time"
)

type Repository struct {
	db      *sql.DB
	ctx     context.Context
	queries *sqlc.Queries
}

func NewRepository(ctx context.Context, db *sql.DB) (*Repository, error) {

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := sqlc.New(db)

	return &Repository{ctx: ctx, db: db, queries: queries}, nil

}

func (r *Repository) GetMealPlanByDate(userID int64, date time.Time) ([]MealPlan, error) {
	monday := getMondayOfCurrentWeek(date)
	params := sqlc.ListMealsParams{
		UserID:   userID,
		StartDay: monday,
		EndDay:   monday.AddDate(0, 0, 7),
	}
	dbMeals, err := r.queries.ListMeals(r.ctx, params)
	if err != nil {
		return nil, err
	}
	if len(dbMeals) == 0 {
		var emptyMeals = make([]MealPlan, len(days))
		for i, day := range days {
			emptyMeals[i] = MealPlan{Day: day, Date: monday.AddDate(0, 0, 1)}
		}
	}

	meals := make([]MealPlan, 7)
	for i, meal := range dbMeals {
		meals[i] = MealPlan{
			Day:       days[i],
			Date:      monday.AddDate(0, 0, i),
			Breakfast: meal.Breakfast,
			Snack1:    meal.Snack1,
			Snack2:    meal.Snack2,
			Lunch:     meal.Lunch,
			Dinner:    meal.Dinner,
		}
	}

	return meals, nil
}

func (r *Repository) UpdateMealsForDate(userID int64, date time.Time, meals []MealPlan) error {
	monday := getMondayOfCurrentWeek(date)
	for i, meal := range meals {
		day := monday.AddDate(0, 0, i)
		params := sqlc.UpdateMealsParams{
			Breakfast: meal.Breakfast,
			Snack1:    meal.Snack1,
			Snack2:    meal.Snack2,
			Lunch:     meal.Lunch,
			Dinner:    meal.Dinner,
			UserID:    userID,
			Day:       day,
		}
		value, err := r.queries.UpdateMeals(r.ctx, params)
		if err != nil {
			return err
		}
		log.Println(value)
	}
	return nil
}
