package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type urlGetter interface {
	GetOriginalURL(ctx context.Context, shortUrl string) (string, error)
}

type redirectService struct {
	logger *zap.Logger
	db     urlGetter
}

func NewRedirectService(logger *zap.Logger, db urlGetter) redirectService {
	return redirectService{
		logger: logger,
		db:     db,
	}
}

func (ss redirectService) FindURL(ctx context.Context, key string) (string, error) {

	baseURL, err := ss.db.GetOriginalURL(ctx, key)

	if err != nil {
		ss.logger.Error(fmt.Sprintf("error while gettingbase url from %s :%s", key, err.Error()))
		return "", err
	}

	return baseURL, nil
}
