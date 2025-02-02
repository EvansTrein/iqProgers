package services

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/storages"
)

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
