package services

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/storages"
)

func (w *Wallet) Deposit(ctx context.Context, req *models.DepositRequest) (*models.DepositResponse, error) {
	op := "service Wallet: deposit request received"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Deposit func call", slog.Any("requets data", req))

	exsistTransaction, err := w.db.ExsistIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		log.Error("failed to check if the transaction exists in the database", "error", err)
		return nil, err
	}

	if exsistTransaction {
		// TODO: транзакция существует, вернуть ее
		return nil, nil
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

	resp := models.DepositResponse{
		Message:   "deposit successfully",
		Operation: dataTran,
	}

	log.Info("deposit successfully")
	return &resp, nil
}
