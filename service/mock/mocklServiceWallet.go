package mock

import (
	"context"

	"github.com/EvansTrein/iqProgers/models"
)

const IdempotencyKeyTestDef = "42dd3893-9baf-43ac-8c2b-32231f486b87"

type MockStoreWallet struct {
	ExsistUserFunc           func(ctx context.Context, id uint) (bool, error)
	ExsistIdempotencyKeyFunc func(ctx context.Context, uuid string) (bool, error)
	TransactionCreateFunc    func(ctx context.Context, data *models.Transaction) error
	TransactionGetFunc       func(ctx context.Context, idempotencyKey string) (*models.Transaction, error)
	DepositFunc              func(ctx context.Context, req *models.DepositRequest) error
	TransferFunc             func(ctx context.Context, req *models.Transaction) error
	OperationsGetFunc        func(ctx context.Context, req *models.UserOperationsRequest) (*models.UserOperationsResponse, error)
}

func (m *MockStoreWallet) ExsistUser(ctx context.Context, id uint) (bool, error) {
	return m.ExsistUserFunc(ctx, id)
}

func (m *MockStoreWallet) ExsistIdempotencyKey(ctx context.Context, uuid string) (bool, error) {
	return m.ExsistIdempotencyKeyFunc(ctx, uuid)
}

func (m *MockStoreWallet) TransactionCreate(ctx context.Context, data *models.Transaction) error {
	return m.TransactionCreateFunc(ctx, data)
}

func (m *MockStoreWallet) TransactionGet(ctx context.Context, idempotencyKey string) (*models.Transaction, error) {
	return m.TransactionGetFunc(ctx, idempotencyKey)
}

func (m *MockStoreWallet) Deposit(ctx context.Context, req *models.DepositRequest) error {
	return m.DepositFunc(ctx, req)
}

func (m *MockStoreWallet) Transfer(ctx context.Context, req *models.Transaction) error {
	return m.TransferFunc(ctx, req)
}

func (m *MockStoreWallet) OperationsGet(ctx context.Context, req *models.UserOperationsRequest) (*models.UserOperationsResponse, error) {
	return m.OperationsGetFunc(ctx, req)
}
