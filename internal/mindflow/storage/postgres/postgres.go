package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgres.New"

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		conn: conn,
	}, nil
}

func (s *Storage) Close() error {
	return s.conn.Close(context.Background())
}
