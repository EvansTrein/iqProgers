package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/EvansTrein/iqProgers/internal/storages"
	"github.com/gin-gonic/gin"
)

// example request

// path parameters - required
// id 1
//
// query parameters 
// limit - required
// 10
//
// offset 
// default = 0
type walletOperations interface {
	UserOperations(ctx context.Context, req *models.UserOperationsRequest) (*models.UserOperationsResponse, error)
}

func Operations(log *slog.Logger, service walletOperations) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Operations: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("request received")

		userID, ok := ctx.Params.Get("id")
		if !ok {
			ctx.JSON(400, models.HandlerResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid data in params",
				Error:   "user id not passed",
			})
			return
		}

		limit, ok := ctx.GetQuery("limit")
		if !ok {
			ctx.JSON(400, models.HandlerResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid data in params",
				Error:   "limit not passed",
			})
			return
		}

		var reqData models.UserOperationsRequest
		params := map[string]string{}

		offset, ok := ctx.GetQuery("offset")
		if !ok {
			log.Debug("offset was not passed, the default value of 0 will be used")
			params["offset"] = "0"
		} else {
			params["offset"] = offset
		}

		params["limit"] = limit
		params["userID"] = userID

		if err := validateRequestParams(params, &reqData); err != nil {
			log.Error("validation params failed")
			ctx.JSON(400, models.HandlerResponse{
				Status:  http.StatusBadRequest,
				Message: "validation params failed",
				Error:   err.Error(),
			})
			return
		}

		log.Debug("request data has been successfully validated", "reqData", reqData)

		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeoutHandlerResponce)
		defer cancel()

		result, err := service.UserOperations(timeoutCtx, &reqData)
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
			case errors.Is(err, storages.ErrOperationsNotFound):
				log.Warn("deposit failed, user has no operations", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Message: "user has no operations",
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

		log.Info("operations successfully")
		ctx.JSON(200, result)
	}
}
