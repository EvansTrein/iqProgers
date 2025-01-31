package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func New(storagePath string, log *slog.Logger) (*PostgresDB, error) {
	log.Debug("database: connection to Postgres started")

	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("database: connect to Postgres successfully")
	return &PostgresDB{db: db, log: log}, nil
}

func (s *PostgresDB) Close() error {
	s.log.Debug("database: stop started")

	if s.db == nil {
		return fmt.Errorf("database connection is already closed")
	}

	s.db.Close()

	s.db = nil

	s.log.Info("database: stop successful")
	return nil
}
