package postgres

import (
	"context"
	"log/slog"
)

func (s *PostgresDB) ExsistIdempotencyKey(ctx context.Context, uuid string) (bool, error) {
	op := "Database: Idempotency Key check"
	log := s.log.With(slog.String("operation", op))
	log.Debug("ExsistIdempotencyKey func call", "uuid", uuid)

	checkQuery := `SELECT EXISTS(SELECT 1 FROM transactions WHERE idempotency_key = $1);`

	var exsist bool
	row := s.db.QueryRow(ctx, checkQuery, uuid)
    if err := row.Scan(&exsist); err != nil {
        log.Error("failed to retrieve data from the database", slog.Any("error", err))
        return false, err
    }

	log.Info("Idempotency check of the key was successful")
	return exsist, nil
}
