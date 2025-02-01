package postgres

import (
	"context"
	"log/slog"
)

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
