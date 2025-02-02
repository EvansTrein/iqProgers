package postgres

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/storages"
)

// OperationsGet retrieves a list of transactions (operations) for a specific user from the database. It queries the database
// for transactions where the user is either the sender or the receiver. The results are ordered by the transaction date in
// descending order and paginated using the provided limit and offset. The function returns a response containing the list
// of transactions, including details such as transaction type, amount, date, and associated sender/receiver names (if applicable).
// If no transactions are found, it returns an error indicating that no operations were found.
// Errors during database querying or row scanning are logged and returned.
func (s *PostgresDB) OperationsGet(ctx context.Context, req *models.UserOperationsRequest) (*models.UserOperationsResponse, error) {
	op := "Database: get user operations "
	log := s.log.With(slog.String("operation", op))
	log.Debug("OperationsGet func call", "data", req)

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
			t.sender_id = $1 OR t.receiver_id = $1
		ORDER BY
			t.date_operation DESC
		LIMIT $2 OFFSET $3;`

	rows, err := s.db.Query(ctx, queryGet, req.UserID, req.Limit, req.Offset)
	if err != nil {
		log.Error("failed to retrieve records from the database", "error", err)
		return nil, err
	}
	defer rows.Close()

	var resp models.UserOperationsResponse
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID,
			&t.Success,
			&t.TypeOperation,
			&t.Amount,
			&t.Date,
			&t.SenderName,
			&t.ReceiverName,
		)
		if err != nil {
			log.Error("failed to scan transaction", "error", err)
			return nil, err
		}
		resp.Operation = append(resp.Operation, &t)
	}

	if err := rows.Err(); err != nil {
		log.Error("error after scanning rows", "error", err)
		return nil, err
	}

	if len(resp.Operation) == 0 {
		log.Warn("user has no operations")
		return nil, storages.ErrOperationsNotFound
	}

	log.Info("transactions were successfully retrieved from the database")
	return &resp, nil
}
