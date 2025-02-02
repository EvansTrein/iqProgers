package postgres

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
	services "github.com/EvansTrein/iqProgers/service"
)

func (s *PostgresDB) Transfer(ctx context.Context, data *models.Transaction) error {
	op := "Database: account transfer"
	log := s.log.With(slog.String("operation", op))
	log.Debug("Transfer func call", "data", data)

	queryLock := `SELECT id FROM users WHERE id = $1 FOR UPDATE;`

	queryCheckBalance := `
	WITH sender_balance AS (
			SELECT balance FROM users WHERE id = $1
		)
	SELECT (balance >= (ROUND($2::numeric, 2) * 100)::bigint)
	FROM sender_balance;`

	queryUpdateSender := `UPDATE users
		SET balance = balance - (ROUND($1::numeric, 2) * 100)::bigint
		WHERE id = $2;`

	queryUpdateReceiver := `UPDATE users
		SET balance = balance + (ROUND($1::numeric, 2) * 100)::bigint
		WHERE id = $2;`
	
	queryCheckNegativeBalance := `SELECT balance >= 0
		FROM users WHERE id = $1;`

	queryGetName := `
	SELECT
		u_sender.name AS sender_name,
		u_receiver.name AS receiver_name
	FROM
		users u_sender
	JOIN
		users u_receiver ON u_receiver.id = $2
	WHERE
		u_sender.id = $1;`

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", "error", err)
		return err
	}

	if _, err := tx.Exec(ctx, queryLock, data.SenderID); err != nil {
		log.Error("failed to execute SQL query lock sender in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	if _, err := tx.Exec(ctx, queryLock, data.ReceiverID); err != nil {
		log.Error("failed to execute SQL query lock receiver in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	var checkBalance bool
	row := tx.QueryRow(ctx, queryCheckBalance, data.SenderID, data.Amount)
	if err != row.Scan(&checkBalance) {
		log.Error("failed to execute SQL query check balance in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	if !checkBalance {
		log.Warn("insufficient account balance")
		tx.Rollback(ctx)
		return services.ErrInsufficientFunds
	}

	if _, err := tx.Exec(ctx, queryUpdateSender, data.Amount, data.SenderID); err != nil {
		log.Error("failed to execute SQL query update sender in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	var isBalanceNonNegative bool
	row = tx.QueryRow(ctx, queryCheckNegativeBalance, data.SenderID)
	if err != row.Scan(&isBalanceNonNegative) {
		log.Error("failed to execute SQL query check negative balance in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	if !isBalanceNonNegative {
		log.Warn("negative balance")
		tx.Rollback(ctx)
		return services.ErrNegaticeBalance
	}

	if _, err := tx.Exec(ctx, queryUpdateReceiver, data.Amount, data.ReceiverID); err != nil {
		log.Error("failed to execute SQL query update receiver in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	row = tx.QueryRow(ctx, queryGetName, data.SenderID, data.ReceiverID)
	if err := row.Scan(&data.SenderName, &data.ReceiverName); err != nil {
		log.Error("failed to execute SQL query get name in the database", "error", err)
		tx.Rollback(ctx)
		return err
	}

	if err := s.TransactionSetResult(ctx, data.IdempotencyKey, true); err != nil {
		log.Error("failed to set the result of user transaction", "error", err)
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("!!!ATTENTION!!! failed to commit transaction", "error", err)
		return err
	}

	return nil
}
