package postgres

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
)

func (s *PostgresDB) Deposit(ctx context.Context, req *models.DepositRequest) error {
	op := "Database: account deposit"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Deposit func call", "data", req)

	queryLock := `SELECT id FROM users WHERE id = $1 FOR UPDATE;`

	updateQuery := `UPDATE users
		SET balance = balance + (ROUND($1::numeric, 2) * 100)::bigint
		WHERE id = $2;`

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", "error", err)
		return err
	}

	if _, err := tx.Exec(ctx, queryLock, req.UserID); err != nil {
		log.Error("failed to execute SQL query lock in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	if _, err := tx.Exec(ctx, updateQuery, req.Amount, req.UserID); err != nil {
		log.Error("failed to execute SQL query to update the balance in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	if err := s.TransactionSetResult(ctx, req.IdempotencyKey, true); err != nil {
		log.Error("failed to set the result of user transaction", "error", err)
		tx.Rollback(ctx)
		return err
	}
	
	if err := tx.Commit(ctx); err != nil {
		log.Error("!!!ATTENTION!!! failed to commit transaction", "error", err)
		return err
	}

	log.Info("transaction successfully completed")
	return nil
}
