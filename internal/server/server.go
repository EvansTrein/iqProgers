package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/EvansTrein/iqProgers/internal/config"
	"github.com/gin-gonic/gin"
)

const gracefulShutdownTimer = time.Second * 20

type HttpServer struct {
	router *gin.Engine
	server *http.Server
	log    *slog.Logger
	conf   *config.HTTPServer
}

func New(log *slog.Logger, conf *config.HTTPServer) *HttpServer {
	router := gin.Default()

	return &HttpServer{
		router: router,
		conf:   conf,
		log:    log,
	}
}

func (s *HttpServer) Start() error {
	log := s.log.With(slog.String("Address", s.conf.Address+":"+s.conf.Port))

	log.Debug("HTTP server: started creating")

	s.server = &http.Server{
		Addr:    s.conf.Address + ":" + s.conf.Port,
		Handler: s.router,
	}

	log.Info("HTTP server: successfully started")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Stop() error {
	s.log.Debug("HTTP server: stop started")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimer)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Server shutdown failed", "error", err)
		return err
	}

	s.server = nil

	s.log.Info("HTTP server: stop successful")
	return nil
}
