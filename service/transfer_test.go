package services

import (
	"context"
	"errors"
	"testing"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/pkg/logs"
	"github.com/EvansTrein/iqProgers/service/mock"
	"github.com/EvansTrein/iqProgers/storages"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWallet_Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logs.NewDiscardLogger()
	mockStore := &mock.MockStoreWallet{}

	wallet := New(log, mockStore)

	tests := []struct {
		name         string
		req          *models.TransferRequest
		mockSetup    func()
		expectedResp *models.TransferResponse
		expectedErr  error
	}{
		{
			name: "successful transfer",
			req: &models.TransferRequest{
				SenderID:       1,
				ReceiverID:     2,
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
				mockStore.TransferFunc = func(ctx context.Context, req *models.Transaction) error {
					return nil
				}
			},
			expectedResp: &models.TransferResponse{
				Message: "transfer successfully",
				Operation: &models.Transaction{
					IdempotencyKey: mock.IdempotencyKeyTestDef,
					SenderID:       1,
					ReceiverID:     2,
					TypeOperation:  "transfer",
					Amount:         100,
					Success:        true,
				},
			},
			expectedErr: nil,
		},
		{
			name: "transaction already exists",
			req: &models.TransferRequest{
				SenderID:       1,
				ReceiverID:     2,
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
						ReceiverID:     2,
						TypeOperation:  "transfer",
						Amount:         100,
						Success:        true,
					}, nil
				}
			},
			expectedResp: &models.TransferResponse{
				Message: "transfer successfully",
				Operation: &models.Transaction{
					IdempotencyKey: mock.IdempotencyKeyTestDef,
					SenderID:       1,
					ReceiverID:     2,
					TypeOperation:  "transfer",
					Amount:         100,
					Success:        true,
				},
			},
			expectedErr: nil,
		},
		{
			name: "user not found",
			req: &models.TransferRequest{
				SenderID:       1,
				ReceiverID:     2,
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
			req: &models.TransferRequest{
				SenderID:       1,
				ReceiverID:     2,
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
		{
			name: "failed to transfer",
			req: &models.TransferRequest{
				SenderID:       1,
				ReceiverID:     2,
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
				mockStore.TransferFunc = func(ctx context.Context, req *models.Transaction) error {
					return errors.New("failed to transfer")
				}
			},
			expectedResp: nil,
			expectedErr:  errors.New("failed to transfer"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := wallet.Transfer(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}
