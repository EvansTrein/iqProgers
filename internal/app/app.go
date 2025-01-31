package app

import (
	"log/slog"

	"github.com/EvansTrein/iqProgers/internal/config"
	"github.com/EvansTrein/iqProgers/internal/server"
	services "github.com/EvansTrein/iqProgers/service"
	"github.com/EvansTrein/iqProgers/storages/postgres"
)

type App struct {
	server *server.HttpServer
	log    *slog.Logger
	conf   *config.Config
	db     *postgres.PostgresDB
	wallet *services.Wallet
}

func New(conf *config.Config, log *slog.Logger) *App {
	log.Debug("application: creation is started")

	httpServer := server.New(log, &conf.HTTPServer)

	db, err := postgres.New(conf.StoragePath, log)
	if err != nil {
		panic(err)
	}

	wallet := services.New(log, db)

	httpServer.InitRouters(wallet)

	return &App{
		server: httpServer,
		log:  log,
		conf: conf,
		db:   db,
		wallet: wallet,
	}
}

func (a *App) MustStart() {
	a.log.Debug("application: started")

	a.log.Info("application: successfully started", "port", a.conf.HTTPServer.Port)
	if err := a.server.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() error {
	a.log.Debug("application: stop started")

	if err := a.server.Stop(); err != nil {
		a.log.Error("failed to stop HTTP server")
		return err
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("failed to close the database connection")
		return err
	}

	if err := a.wallet.Stop(); err != nil {
		a.log.Error("failed to stop the Wallet service")
		return err
	}

	a.server = nil
	a.wallet = nil
	a.db = nil

	a.log.Info("application: stop successful")
	return nil
}