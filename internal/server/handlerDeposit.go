package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/storages"
	"github.com/EvansTrein/iqProgers/utils"
	"github.com/gin-gonic/gin"
)

type walletDeposit interface {
	Deposit(ctx context.Context, req *models.DepositRequest) (*models.DepositResponse, error)
}

// example request
//
// Headers - required
// Idempotency-Key UUID
// 'f65616ca-8b51-4af2-8342-84157b55cbb7'
//
// body - required
// {
// "id": 2, 
// "amount": 205.44
// }
func Deposit(log *slog.Logger, service walletDeposit) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Deposit: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("request received")

		var reqData models.DepositRequest
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

		if checkFormat := utils.IsGUID(reqData.IdempotencyKey); !checkFormat {
			ctx.JSON(400, models.HandlerResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid data in headers",
				Error:   "'Idempotency-Key' of incorrect format",
			})
			return
		}

		log.Debug("request data has been successfully validated", "reqData", reqData)

		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeoutHandlerResponce)
		defer cancel()

		result, err := service.Deposit(timeoutCtx, &reqData)
		if err != nil {
			switch {
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

		log.Info("deposit successfully")
		ctx.JSON(200, result)
	}
}
