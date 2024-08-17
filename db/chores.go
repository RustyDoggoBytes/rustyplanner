package db

import (
	"database/sql"
	sqlc "rustydoggobytes/planner/sqlc_generated"
	"time"
)

type Chores struct {
	ID       int64
	Title    string
	Assigned string
	DueDate  time.Time
}

func (r *Repository) GetChores(userID int64) ([]Chores, error) {
	dbChores, err := r.queries.ListOnceChores(r.ctx, userID)

	if err != nil {
		return nil, err
	}

	chores := make([]Chores, len(dbChores))
	for i, db := range dbChores {
		chores[i] = mapOnceChore(db.Chore, db.ChoresRecurrenceOnce)
	}

	return chores, nil
}

func mapOnceChore(chore sqlc.Chore, onceRecurrence sqlc.ChoresRecurrenceOnce) Chores {
	assigned := ""
	if chore.Assigned.Valid {
		assigned = chore.Assigned.String
	}
	return Chores{
		ID:       chore.ID,
		Title:    chore.Title,
		Assigned: assigned,
		DueDate:  onceRecurrence.DueDate,
	}
}

func (r *Repository) CreateChore(userID int64, title, assigned string, dueDate time.Time) (*Chores, error) {
	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	recurrence, err := qtx.CreateOnceRecurrence(r.ctx, dueDate)
	if err != nil {
		return nil, err
	}

	dbChore, err := qtx.CreateChore(r.ctx, sqlc.CreateChoreParams{
		Title:          title,
		UserID:         userID,
		RecurrenceID:   recurrence.ID,
		RecurrenceType: "once",
		Assigned:       sql.NullString{String: assigned, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	chore := mapOnceChore(dbChore, recurrence)
	return &chore, nil
}

func (r *Repository) DeleteChore(userID int64, choreID int64) error {
	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := r.queries.WithTx(tx)

	deletedRow, err := qtx.DeleteChore(r.ctx, sqlc.DeleteChoreParams{
		ID:     choreID,
		UserID: userID,
	})
	if err != nil {
		return err
	}
	_, err = qtx.DeleteChoreRecurrenceOnce(r.ctx, deletedRow.RecurrenceID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
