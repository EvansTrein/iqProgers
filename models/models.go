package models

import "time"

type HandlerResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type DepositRequest struct {
	IdempotencyKey string  `json:"-"`
	UserID         uint    `json:"id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
}

type DepositResponse struct {
	Message   string       `json:"message"`
	Operation *Transaction `json:"operation"`
}

type TransferRequest struct {
	IdempotencyKey string  `json:"-"`
	SenderID       uint    `json:"sender_id" binding:"required"`
	ReceiverID     uint    `json:"receiver_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
}

type TransferResponse struct {
	Message   string       `json:"message"`
	Operation *Transaction `json:"operation"`
}

type UserOperationsRequest struct {
	UserID uint
	Offset int
	Limit  int
}

type UserOperationsResponse struct {
	Message   string         `json:"message"`
	Operation []*Transaction `json:"operation"`
}

type Transaction struct {
	ID             uint      `json:"transaction_id"`
	SenderID       uint      `json:"-"`
	ReceiverID     uint      `json:"-"`
	IdempotencyKey string    `json:"-"`
	Success        bool      `json:"success"`
	SenderName     *string   `json:"sender,omitempty"`
	ReceiverName   *string   `json:"receiver,omitempty"`
	TypeOperation  string    `json:"type_operation"`
	Amount         float64   `json:"amount"`
	Date           time.Time `json:"date"`
}

type User struct {
	ID      uint
	Name    string
	Balance float64
}
