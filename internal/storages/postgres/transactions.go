package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/EvansTrein/iqProgers/models"
)

// TransactionCreate inserts a new transaction record into the database. Depending on the transaction type (deposit or transfer),
// it executes the appropriate SQL query to create the transaction. For deposits, it inserts the sender ID, idempotency key,
// transaction type, and amount. For transfers, it additionally includes the receiver ID. The function retrieves the newly
// created transaction's ID and date of operation, updates the provided transaction data with these values, and returns nil
// on success. If the database operation fails, the error is logged and returned. This function ensures that transactions
// are recorded accurately and consistently in the database.
func (s *PostgresDB) TransactionCreate(ctx context.Context, data *models.Transaction) error {
	op := "Database: transaction creation"
	log := s.log.With(slog.String("operation", op))
	log.Debug("createTransaction func call", "data", data)

	createDepositQuery := `INSERT INTO transactions
		(sender_id, idempotency_key, type_operation, amount)
		VALUES
		($1, $2, $3, (ROUND($4::numeric, 2) * 100)::bigint)
		RETURNING id, date_operation;`

	createTransferQuery := `INSERT INTO transactions
		(sender_id, receiver_id, idempotency_key, type_operation, amount)
		VALUES
		($1, $2, $3, $4, (ROUND($5::numeric, 2) * 100)::bigint)
		RETURNING id, date_operation;`

	var id uint
	var dateOperation time.Time

	switch data.TypeOperation {
	case "deposit":
		row := s.db.QueryRow(ctx, createDepositQuery, data.SenderID, data.IdempotencyKey, data.TypeOperation, data.Amount)
		if err := row.Scan(&id, &dateOperation); err != nil {
			log.Error("failed to create transaction", "error", err)
			return err
		}
	case "transfer":
		row := s.db.QueryRow(ctx, createTransferQuery, data.SenderID, data.ReceiverID, data.IdempotencyKey, data.TypeOperation, data.Amount)
		if err := row.Scan(&id, &dateOperation); err != nil {
			log.Error("failed to create transaction", "error", err)
			return err
		}
	}

	data.ID = id
	data.Date = dateOperation

	log.Info("transaction created successfully", "id", id)
	return nil
}

// TransactionSetResult updates the success status of a transaction in the database using the provided idempotency key.
// It executes an SQL query to set the `success` field of the transaction record to the specified value (true or false).
// If the update operation fails, the error is logged and returned. This function is used to mark the outcome of a transaction
func (s *PostgresDB) TransactionSetResult(ctx context.Context, idempotencyKey string, success bool) error {
	op := "Database: transaction result"
	log := s.log.With(slog.String("operation", op))
	log.Debug("TransactionSetResult func call", "success", success)

	resultQuery := `UPDATE transactions
		SET success = $1
		WHERE idempotency_key = $2;`

	if _, err := s.db.Exec(ctx, resultQuery, success, idempotencyKey); err != nil {
		log.Error("failed to update the user transaction result in the database", "error", err)
		return err
	}

	log.Info("user transaction result in the database was successfully updated")
	return nil
}

// TransactionGet retrieves a transaction from the database using the provided idempotency key. It queries the database
// to fetch details of the transaction, including its ID, success status, type, amount, date, and associated sender/receiver
// names (if applicable). The function joins the `transactions` table with the `users` table to retrieve sender and receiver
// names for non-deposit transactions. If the query fails or the transaction is not found, the error is logged and returned.
func (s *PostgresDB) TransactionGet(ctx context.Context, idempotencyKey string) (*models.Transaction, error) {
	op := "Database: get transactions"
	log := s.log.With(slog.String("operation", op))
	log.Debug("TransactionGet func call", "idempotencyKey", idempotencyKey)

	queryGet := `
		SELECT
			t.id,
			t.success,
			t.type_operation,
			CAST(t.amount AS FLOAT) / 100 AS amount,
			t.date_operation,
			CASE
				WHEN t.type_operation != 'deposit' THEN u_sender.name
				ELSE NULL
			END AS sender_name,
			CASE
				WHEN t.type_operation != 'deposit' THEN u_receiver.name
				ELSE NULL
			END AS receiver_name
		FROM
			transactions t
		LEFT JOIN
			users u_sender ON t.sender_id = u_sender.id
		LEFT JOIN
			users u_receiver ON t.receiver_id = u_receiver.id
		WHERE
			t.idempotency_key = $1;`

	var transaction models.Transaction

	row := s.db.QueryRow(ctx, queryGet, idempotencyKey)
	if err := row.Scan(
		&transaction.ID,
		&transaction.Success,
		&transaction.TypeOperation,
		&transaction.Amount,
		&transaction.Date,
		&transaction.SenderName,
		&transaction.ReceiverName,
	); err != nil {
		log.Error("failed to get the transaction", "error", err)
		return nil, err
	}

	log.Debug("data was retrieved from the database", "transaction", transaction)

	log.Info("transaction is successfully retrieved from the database")
	return &transaction, nil
}
