package db

import (
	"context"
	"database/sql"
	sqlc "rustydoggobytes/planner/sqlc_generated"
)

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
