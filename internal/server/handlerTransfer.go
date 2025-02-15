package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/EvansTrein/iqProgers/models"
	serv "github.com/EvansTrein/iqProgers/internal/service"
	"github.com/EvansTrein/iqProgers/internal/storages"
	"github.com/gin-gonic/gin"
)

type walletTransfer interface {
	Transfer(ctx context.Context, req *models.TransferRequest) (*models.TransferResponse, error)
}

// example request
//
// Headers - required
// Idempotency-Key UUID
// 'f65616ca-8b51-4af2-8342-84157b55cbb7'
//
// body - required
//
//	{
//		"sender_id": 4,
//		"receiver_id": 3,
//		"amount": 100.55
//	}
func Transfer(log *slog.Logger, service walletTransfer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Transfer: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("request received")

		var reqData models.TransferRequest
		if err := ctx.ShouldBindJSON(&reqData); err != nil {
			ctx.JSON(400, models.HandlerResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid data in body",
				Error:   err.Error(),
			})
			return
		}

		reqData.IdempotencyKey = ctx.GetHeader("Idempotency-Key")
		if reqData.IdempotencyKey == "" {
			ctx.JSON(400, models.HandlerResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid data in headers",
				Error:   "'Idempotency-Key' was not passed in the headers",
			})
			return
		}

		isVaild, err := isGUID(reqData.IdempotencyKey)
		if err != nil {
			ctx.JSON(500, models.HandlerResponse{
				Status:  http.StatusInternalServerError,
				Message: "failed to verify the header format 'Idempotency-Key'",
				Error:   err.Error(),
			})
			return
		}

		if !isVaild {
			ctx.JSON(400, models.HandlerResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid data in headers 'Idempotency-Key'",
				Error:   "header 'Idempotency-Key' does not match the UUID format",
			})
			return
		}

		log.Debug("request data has been successfully validated", "reqData", reqData)

		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeoutHandlerResponce)
		defer cancel()

		result, err := service.Transfer(timeoutCtx, &reqData)
		if err != nil {
			switch {
			case errors.Is(err, serv.ErrInsufficientFunds):
				log.Warn("transfer failed, insufficient funds on the balance sheet", "error", err)
				ctx.JSON(402, models.HandlerResponse{
					Status:  http.StatusPaymentRequired,
					Message: "insufficient funds",
					Error:   err.Error(),
				})
				return
			case errors.Is(err, serv.ErrNegaticeBalance):
				log.Error("transfer failed, negative balance", "error", err)
				ctx.JSON(422, models.HandlerResponse{
					Status:  http.StatusUnprocessableEntity,
					Message: "balance cannot be negative",
					Error:   err.Error(),
				})
				return
			case errors.Is(err, storages.ErrUserNotFound):
				log.Warn("deposit failed, no user with this id", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Message: "no user with this id",
					Error:   err.Error(),
				})
				return
			case errors.Is(err, context.DeadlineExceeded):
				log.Error("deposit failed due to timeout", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Message: "deposit failed due to timeout",
					Error:   err.Error(),
				})
				return
			default:
				log.Error("deposit failed", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status:  http.StatusInternalServerError,
					Message: "deposit failed",
					Error:   err.Error(),
				})
				return
			}
		}

		log.Info("transfer successfully")
		ctx.JSON(200, result)
	}
}
