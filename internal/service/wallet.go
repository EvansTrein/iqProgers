package services

import (
	"errors"
	"log/slog"

	"github.com/EvansTrein/iqProgers/internal/storages"
)

var (
	ErrInsufficientFunds = errors.New("insufficient account balance")
	ErrNegaticeBalance = errors.New("negative balance")
)

type Wallet struct {
	log *slog.Logger
	db  storages.StoreWallet
}

func New(log *slog.Logger, db storages.StoreWallet) *Wallet {
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
