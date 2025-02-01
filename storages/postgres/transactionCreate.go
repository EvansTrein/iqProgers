package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/EvansTrein/iqProgers/models"
)

func (s *PostgresDB) TransactionCreate(ctx context.Context, data *models.Transaction) error {
	op := "Database: transaction creation"
	log := s.log.With(slog.String("operation", op))
	log.Debug("createTransaction func call", "data", data)

	createQuery := `INSERT INTO transactions
		(sender_id, idempotency_key, type_operation, amount)
		VALUES
		($1, $2, $3, (ROUND($4::numeric, 2) * 100)::bigint)
		RETURNING id, date_operation;`

	var id uint
	var dateOperation time.Time
	row := s.db.QueryRow(ctx, createQuery, data.SenderID, data.IdempotencyKey, data.TypeOperation, data.Amount)
	if err := row.Scan(&id, &dateOperation); err != nil {
		log.Error("failed to create transaction", "error", err)
		return err
	}

	data.ID = id
	data.Date = dateOperation

	log.Info("transaction created successfully", "id", id)
	return nil
}
