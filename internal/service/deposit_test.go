package services

import (
	"context"
	"errors"
	"testing"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/pkg/logs"
	"github.com/EvansTrein/iqProgers/internal/service/mock"
	"github.com/EvansTrein/iqProgers/internal/storages"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWallet_Deposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logs.NewDiscardLogger()
	mockStore := &mock.MockStoreWallet{}

	wallet := New(log, mockStore)

	tests := []struct {
		name         string
		req          *models.DepositRequest
		mockSetup    func()
		expectedResp *models.DepositResponse
		expectedErr  error
	}{
		{
			name: "successful deposit",
			req: &models.DepositRequest{
				UserID:         1,
				Amount:         100,
				IdempotencyKey: mock.IdempotencyKeyTestDef,
			},
			mockSetup: func() {
				mockStore.ExsistIdempotencyKeyFunc = func(ctx context.Context, uuid string) (bool, error) {
					return false, nil
				}
				mockStore.ExsistUserFunc = func(ctx context.Context, id uint) (bool, error) {
					return true, nil
				}
				mockStore.TransactionCreateFunc = func(ctx context.Context, data *models.Transaction) error {
					return nil
				}
				mockStore.DepositFunc = func(ctx context.Context, req *models.DepositRequest) error {
					return nil
				}
			},
			expectedResp: &models.DepositResponse{
				Message: "deposit successfully",
				Operation: &models.Transaction{
					IdempotencyKey: mock.IdempotencyKeyTestDef,
					SenderID:       1,
					TypeOperation:  "deposit",
					Amount:         100,
					Success:        true,
				},
			},
			expectedErr: nil,
		},
		{
			name: "transaction already exists",
			req: &models.DepositRequest{
				UserID:         1,
				Amount:         100,
				IdempotencyKey: mock.IdempotencyKeyTestDef,
			},
			mockSetup: func() {
				mockStore.ExsistIdempotencyKeyFunc = func(ctx context.Context, uuid string) (bool, error) {
					return true, nil
				}
				mockStore.TransactionGetFunc = func(ctx context.Context, idempotencyKey string) (*models.Transaction, error) {
					return &models.Transaction{
						IdempotencyKey: mock.IdempotencyKeyTestDef,
						SenderID:       1,
						TypeOperation:  "deposit",
						Amount:         100,
						Success:        true,
					}, nil
				}
			},
			expectedResp: &models.DepositResponse{
				Message: "deposit successfully",
				Operation: &models.Transaction{
					IdempotencyKey: mock.IdempotencyKeyTestDef,
					SenderID:       1,
					TypeOperation:  "deposit",
					Amount:         100,
					Success:        true,
				},
			},
			expectedErr: nil,
		},
		{
			name: "user not found",
			req: &models.DepositRequest{
				UserID:         1,
				Amount:         100,
				IdempotencyKey: mock.IdempotencyKeyTestDef,
			},
			mockSetup: func() {
				mockStore.ExsistIdempotencyKeyFunc = func(ctx context.Context, uuid string) (bool, error) {
					return false, nil
				}
				mockStore.ExsistUserFunc = func(ctx context.Context, id uint) (bool, error) {
					return false, nil
				}
			},
			expectedResp: nil,
			expectedErr:  storages.ErrUserNotFound,
		},
		{
			name: "failed to create transaction",
			req: &models.DepositRequest{
				UserID:         1,
				Amount:         100,
				IdempotencyKey: mock.IdempotencyKeyTestDef,
			},
			mockSetup: func() {
				mockStore.ExsistIdempotencyKeyFunc = func(ctx context.Context, uuid string) (bool, error) {
					return false, nil
				}
				mockStore.ExsistUserFunc = func(ctx context.Context, id uint) (bool, error) {
					return true, nil
				}
				mockStore.TransactionCreateFunc = func(ctx context.Context, data *models.Transaction) error {
					return errors.New("failed to create transaction")
				}
			},
			expectedResp: nil,
			expectedErr:  errors.New("failed to create transaction"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := wallet.Deposit(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}
