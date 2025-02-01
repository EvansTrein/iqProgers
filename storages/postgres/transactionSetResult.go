package postgres

import (
	"context"
	"log/slog"
)

func (s *PostgresDB) TransactionSetResult(ctx context.Context, idempotencyKey string, success, completed bool) error {
	op := "Database: transaction result"
	log := s.log.With(slog.String("operation", op))
	log.Debug("TransactionSetResult func call", "success", success, "completed", completed)

	return nil
}
