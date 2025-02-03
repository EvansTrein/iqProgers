package services

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/internal/storages"
)

// UserOperations retrieves the list of operations (transactions) for a specific user.
// It first verifies the existence of the user in the database.
// If the user does not exist, it returns an error indicating the user was not found.
// If the user exists, it fetches the user's transactions
// from the database and returns them in a response. Errors during database access or user verification are logged and returned.
// The response includes a success message and the list of transactions associated with the user.
func (w *Wallet) UserOperations(ctx context.Context, req *models.UserOperationsRequest) (*models.UserOperationsResponse, error) {
	op := "service Wallet: user operations request received"
	log := w.log.With(slog.String("operation", op))
	log.Debug("UserOperations func call", "requets data", req)

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

	resp, err := w.db.OperationsGet(ctx, req)
	if err != nil {
		log.Error("failed to retrieve operations from the database", "error", err)
		return nil, err
	}

	resp.Message = "transactions have been successfully received"

	log.Info("user operations successfully")
	return resp, nil
}
