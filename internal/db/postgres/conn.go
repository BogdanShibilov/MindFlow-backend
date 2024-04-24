package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	*pgxpool.Pool
}

func New(connString string) (*Db, error) {
	const op = "db.postgres.New"

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Db{
		Pool: pool,
	}, nil
}

func (db *Db) Close() {
	db.Pool.Close()
}
