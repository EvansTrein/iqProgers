package postgres

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
)

// Deposit processes a deposit request for a user's account. It locks the user's row in the database to prevent concurrent updates,
// updates the user's balance by adding the specified amount, and marks the transaction as successful. The function uses a database
// transaction to ensure atomicity. If any step fails (e.g., SQL query execution, transaction commit), 
// the transaction is rolled back, and the error is logged and returned. 
func (s *PostgresDB) Deposit(ctx context.Context, req *models.DepositRequest) error {
	op := "Database: account deposit"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Deposit func call", "data", req)

	rollbackCtx := context.Background()

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
		if err := tx.Rollback(rollbackCtx); err != nil {
			log.Error("!!!ATTENTION!!! failed to rollback transaction", "error", err)
		}
		return err
	}

	if _, err := tx.Exec(ctx, updateQuery, req.Amount, req.UserID); err != nil {
		log.Error("failed to execute SQL query to update the balance in the database", "error", err)
		if err := tx.Rollback(rollbackCtx); err != nil {
			log.Error("!!!ATTENTION!!! failed to rollback transaction", "error", err)
		}
		return err
	}

	if err := s.TransactionSetResult(ctx, req.IdempotencyKey, true); err != nil {
		log.Error("failed to set the result of user transaction", "error", err)
		if err := tx.Rollback(rollbackCtx); err != nil {
			log.Error("!!!ATTENTION!!! failed to rollback transaction", "error", err)
		}
		return err
	}
	
	if err := tx.Commit(ctx); err != nil {
		log.Error("!!!ATTENTION!!! failed to commit transaction", "error", err)
		return err
	}

	log.Info("transaction successfully completed")
	return nil
}
