package storages

import (
	"context"
	"errors"

	"github.com/EvansTrein/iqProgers/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	// ErrIdempotencyKeyAlreadyExists = errors.New("Idempotency-Key already exists")
)

type StoreWallet interface {
	ExsistUser(ctx context.Context, id uint) (bool, error)
	ExsistIdempotencyKey(ctx context.Context, uuid string) (bool, error)
	TransactionCreate(ctx context.Context, data *models.Transaction) error
	Deposit(ctx context.Context, req *models.DepositRequest) error
}
