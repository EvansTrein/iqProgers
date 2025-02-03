package services

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/internal/storages"
)

// Deposit handles the deposit request for a user's wallet. It checks if the transaction already exists using the idempotency key.
// If the transaction exists, it retrieves and returns the existing transaction details.
// If the transaction does not exist, it verifies the user's existence, creates a new transaction, updates the user's balance,
// and marks the transaction as successful. The function returns a response indicating the success of the deposit operation.
func (w *Wallet) Deposit(ctx context.Context, req *models.DepositRequest) (*models.DepositResponse, error) {
	op := "service Wallet: deposit request received"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Deposit func call", "requets data", req)

	exsistTransaction, err := w.db.ExsistIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		log.Error("failed to check if the transaction exists in the database", "error", err)
		return nil, err
	}

	if exsistTransaction {
		log.Warn("transaction already exists")
		
		dataTran, err := w.db.TransactionGet(ctx, req.IdempotencyKey)
		if err != nil {
			log.Error("failed to retrieve existing transaction", "error", err)
			return nil, err
		}

		resp := models.DepositResponse{
			Message:   "deposit successfully",
			Operation: dataTran,
		}

		log.Warn("existing transaction successfully sent")
		return &resp, nil
	}

	exsistUser, err := w.db.ExsistUser(ctx, req.UserID)
	if err != nil {
		log.Error("failed to check if the user exists in the database", "error", err)
		return nil, err
	}

	if !exsistUser {
		log.Warn("user not found", "id", req.UserID)
		return nil, storages.ErrUserNotFound
	}

	log.Debug("request data successfully verified")

	// data for transaction creation
	dataTran := models.Transaction{
		IdempotencyKey: req.IdempotencyKey,
		SenderID:       req.UserID,
		TypeOperation:  "deposit",
		Amount:         req.Amount,
	}

	if err := w.db.TransactionCreate(ctx, &dataTran); err != nil {
		log.Error("failed to create a transaction for user operation", "error", err)
		return nil, err
	}

	log.Info("transaction for the user operation was successfully created", "transaction ID", dataTran.ID)

	if err := w.db.Deposit(ctx, req); err != nil {
		log.Error("failed to update the balance value in the database", "error", err)
		return nil, err
	}

	dataTran.Success = true

	resp := models.DepositResponse{
		Message:   "deposit successfully",
		Operation: &dataTran,
	}

	log.Info("deposit successfully")
	return &resp, nil
}
