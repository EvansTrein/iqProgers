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

func TestWallet_UserOperations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logs.NewDiscardLogger()
	mockStore := &mock.MockStoreWallet{}

	wallet := New(log, mockStore)

	tests := []struct {
		name         string
		req          *models.UserOperationsRequest
		mockSetup    func()
		expectedResp *models.UserOperationsResponse
		expectedErr  error
	}{
		{
			name: "successful retrieval of operations",
			req: &models.UserOperationsRequest{
				UserID: 1,
			},
			mockSetup: func() {
				mockStore.ExsistUserFunc = func(ctx context.Context, id uint) (bool, error) {
					return true, nil
				}
				mockStore.OperationsGetFunc = func(ctx context.Context, req *models.UserOperationsRequest) (*models.UserOperationsResponse, error) {
					return &models.UserOperationsResponse{
						Message: "transactions have been successfully received",
						Operation: []*models.Transaction{
							{
								IdempotencyKey: mock.IdempotencyKeyTestDef,
								SenderID:       1,
								TypeOperation:  "deposit",
								Amount:         100,
								Success:        true,
							},
						},
					}, nil
				}
			},
			expectedResp: &models.UserOperationsResponse{
				Message: "transactions have been successfully received",
				Operation: []*models.Transaction{
					{
						IdempotencyKey: mock.IdempotencyKeyTestDef,
						SenderID:       1,
						TypeOperation:  "deposit",
						Amount:         100,
						Success:        true,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "user not found",
			req: &models.UserOperationsRequest{
				UserID: 1,
			},
			mockSetup: func() {
				mockStore.ExsistUserFunc = func(ctx context.Context, id uint) (bool, error) {
					return false, nil
				}
			},
			expectedResp: nil,
			expectedErr:  storages.ErrUserNotFound,
		},
		{
			name: "failed to retrieve operations",
			req: &models.UserOperationsRequest{
				UserID: 1,
			},
			mockSetup: func() {
				mockStore.ExsistUserFunc = func(ctx context.Context, id uint) (bool, error) {
					return true, nil
				}
				mockStore.OperationsGetFunc = func(ctx context.Context, req *models.UserOperationsRequest) (*models.UserOperationsResponse, error) {
					return nil, errors.New("failed to retrieve operations")
				}
			},
			expectedResp: nil,
			expectedErr:  errors.New("failed to retrieve operations"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := wallet.UserOperations(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}
