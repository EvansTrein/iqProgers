package postgres

import (
	"context"
	"log/slog"
)

// ExsistIdempotencyKey checks if a transaction with the given idempotency key already exists in the database.
// It queries the database to determine if a record with the specified idempotency key is present.
// If the query fails or the database returns an error, the function logs the error and returns it.
// Otherwise, it returns a boolean indicating whether the idempotency key exists and a nil error.
// This function is used to ensure idempotency in transaction processing.
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

// ExsistUser checks if a user with the specified ID exists in the database. It queries the database to determine
// if a record with the given user ID is present. If the query fails or the database returns an error, the function
// logs the error and returns it. Otherwise, it returns a boolean indicating whether the user exists and a nil error.
// This function is used to verify the existence of a user before performing operations that require a valid user.
func (s *PostgresDB) ExsistUser(ctx context.Context, id uint) (bool, error) {
	op := "Database: user check"
	log := s.log.With(slog.String("operation", op))
	log.Debug("ExsistUser func call", "user id", id)

	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);`

	var exsist bool
	row := s.db.QueryRow(ctx, checkQuery, id)
	if err := row.Scan(&exsist); err != nil {
		log.Error("failed to retrieve data from the database", slog.Any("error", err))
		return false, err
	}

	log.Info("user checked was successful")
	return exsist, nil
}