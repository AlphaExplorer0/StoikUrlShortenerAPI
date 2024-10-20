package service

import (
	"crypto/md5"
	"math/big"

	"go.uber.org/zap"
)

const shortLinkSize int = 7

type shortenerService struct {
	logger *zap.Logger
}

func NewShortenerService(logger *zap.Logger) shortenerService {
	return shortenerService{
		logger: logger,
	}
}

func (ss shortenerService) ShortenURL(url string) string {

	tiny := hashAndCut(url, shortLinkSize)

	// TODO: check in DB if short url exists. If yes, redo the hash

	return "http://" + tiny
}

func hashAndCut(str string, size int) string {
	hash := md5.Sum([]byte(str))
	encodedHash := toBase62([]byte(hash[:]))
	return string(encodedHash[:size])
}

func toBase62(data []byte) string {
	var i big.Int
	i.SetBytes(data[:])
	return i.Text(62)
}
