package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"
)

const shortLinkSize int = 7

type urlStorage interface {
	PutNewURL(ctx context.Context, originalUrl, shortUrl string) error
}

type shortenerService struct {
	logger *zap.Logger
	db     urlStorage
}

func NewShortenerService(logger *zap.Logger, db urlStorage) shortenerService {
	return shortenerService{
		logger: logger,
		db:     db,
	}
}

func (ss shortenerService) ShortenURL(ctx context.Context, url string) (string, error) {

	tiny := hashAndCut(url, shortLinkSize)

	err := ss.db.PutNewURL(ctx, url, tiny)

	if err != nil {
		ss.logger.Error(fmt.Sprintf("error while inserting new short url :%s", err.Error()))
		return "", err
	}

	return "http://" + tiny, nil
}

func hashAndCut(str string, size int) string {
	str += time.Now().String() // to avoid collisions
	hash := md5.Sum([]byte(str))
	encodedHash := toBase62([]byte(hash[:]))
	return string(encodedHash[:size])
}

func toBase62(data []byte) string {
	var i big.Int
	i.SetBytes(data[:])
	return i.Text(62)
}
