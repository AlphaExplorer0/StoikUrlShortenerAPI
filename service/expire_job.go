package service

import (
	"fmt"

	"github.com/robfig/cron"
	"go.uber.org/zap"
)

type urlStore interface {
	ExpireURLs() error
}

type expirationService struct {
	logger *zap.Logger
	db     urlStore
	cr     *cron.Cron
}

func NewExpirationService(logger *zap.Logger, db urlStore) expirationService {
	return expirationService{
		logger: logger,
		db:     db,
	}
}

func (es *expirationService) Start() {
	c := cron.New()

	c.AddFunc("@every 00h05m00s", func() {
		es.Expire()
	})

	es.cr = c
	es.cr.Start()
	es.logger.Info("expiration job starts")
}

func (es *expirationService) Stop() {
	es.cr.Stop()
}

func (es *expirationService) Expire() {
	err := es.db.ExpireURLs()
	if err != nil {
		es.logger.Error(fmt.Sprintf("failed to execute expiration job : %s", err.Error()))
	} else {
		es.logger.Info("expiration job successful")
	}
}
