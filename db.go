package main

import (
	"context"
	"database/sql"
	"log/slog"
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

func (r *Repository) GetMealPlanByDate(userID int64, startDate time.Time, endDate time.Time) ([]MealPlan, error) {
	params := sqlc.ListMealsParams{
		UserID:   userID,
		StartDay: startDate,
		EndDay:   endDate,
	}
	dbMeals, err := r.queries.ListMeals(r.ctx, params)
	if err != nil {
		return nil, err
	}

	if len(dbMeals) == 0 {
		slog.Info("empty meals", "start", startDate)
		var emptyMeals = make([]MealPlan, len(days))
		for i, day := range days {
			emptyMeals[i] = MealPlan{Day: day, Date: startDate.AddDate(0, 0, i)}
		}
		return emptyMeals, nil
	}

	meals := make([]MealPlan, 7)
	for i, meal := range dbMeals {
		meals[i] = MealPlan{
			Day:       days[i],
			Date:      meal.Day,
			Breakfast: meal.Breakfast,
			Snack1:    meal.Snack1,
			Snack2:    meal.Snack2,
			Lunch:     meal.Lunch,
			Dinner:    meal.Dinner,
		}
	}

	return meals, nil
}

func (r *Repository) UpdateMealsForDate(userID int64, meals []MealPlan) error {
	for _, meal := range meals {
		params := sqlc.UpdateMealsParams{
			UserID:    userID,
			Day:       meal.Date,
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
	}
	return nil
}
