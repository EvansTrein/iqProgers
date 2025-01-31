package services

import (
	"log/slog"

	"github.com/EvansTrein/iqProgers/models"
)

type StoreWallet interface {
	TestDB() (int, error)
}

type Wallet struct {
	log *slog.Logger
	db  StoreWallet
}

func New(log *slog.Logger, db StoreWallet) *Wallet {
	log.Debug("service Wallet: started creating")

	log.Info("service Wallet: successfully created")
	return &Wallet{
		log: log,
		db:  db,
	}
}

func (w *Wallet) Stop() error {
	w.log.Debug("service Wallet: stop started")

	w.db = nil

	w.log.Info("service Wallet: stop successful")
	return nil
}

func (w *Wallet) TestWallet() (*models.PingStruct, error) {
	w.log.Debug("TestWallet serv method")
	res, err := w.db.TestDB()
	if err != nil {
		return nil, err
	}

	return &models.PingStruct{ResultDB: res, ResultServ: "successful wallet service"}, nil
}
