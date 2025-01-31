package server

import (
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
	"github.com/gin-gonic/gin"
)

type walletTest interface {
	TestWallet() (*models.PingStruct, error)
}

func HandlerTest(log *slog.Logger, serv walletTest) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler HadlerTest: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("test request")

		result, err := serv.TestWallet()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
		}

		ctx.JSON(200, result)
	}
}
