package services

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/storages"
)

// Transfer handles the transfer of funds between two users. It first checks if the transaction already exists using the idempotency key.
// If the transaction exists, it retrieves and returns the existing transaction details. If the transaction does not exist, it verifies
// the existence of both the sender and receiver users. If either user is not found, it returns an error. If both users exist,
// it creates a new transaction, processes the transfer, and updates the balances in the database. The function returns a response
// indicating the success of the transfer operation.
func (w *Wallet) Transfer(ctx context.Context, req *models.TransferRequest) (*models.TransferResponse, error) {
	op := "service Wallet: transfer request received"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Transfer func call", "requets data", req)

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

		resp := models.TransferResponse{
			Message:   "transfer successfully",
			Operation: dataTran,
		}

		log.Warn("existing transaction successfully sent")
		return &resp, nil
	}

	exsistUserSender, err := w.db.ExsistUser(ctx, req.SenderID)
	if err != nil {
		log.Error("failed to check if the UserSender exists in the database", "error", err)
		return nil, err
	}

	exsistUserReceiver, err := w.db.ExsistUser(ctx, req.SenderID)
	if err != nil {
		log.Error("failed to check if the UserReceiver exists in the database", "error", err)
		return nil, err
	}

	if !exsistUserSender || !exsistUserReceiver {
		log.Warn("user not found", "SenderID", req.SenderID, "ReceiverID", req.ReceiverID)
		return nil, storages.ErrUserNotFound
	}

	log.Debug("request data successfully verified")

	// data for transaction creation
	dataTran := models.Transaction{
		IdempotencyKey: req.IdempotencyKey,
		SenderID:       req.SenderID,
		ReceiverID:     req.ReceiverID,
		TypeOperation:  "transfer",
		Amount:         req.Amount,
	}

	if err := w.db.TransactionCreate(ctx, &dataTran); err != nil {
		log.Error("failed to create a transaction for user operation", "error", err)
		return nil, err
	}

	log.Info("transaction for the user operation was successfully created", "transaction ID", dataTran.ID)

	if err := w.db.Transfer(ctx, &dataTran); err != nil {
		log.Error("failed to update the balance value in the database", "error", err)
		return nil, err
	}

	dataTran.Success = true

	resp := models.TransferResponse{
		Message:   "transfer successfully",
		Operation: &dataTran,
	}

	log.Info("transfer successfully")
	return &resp, nil
}
