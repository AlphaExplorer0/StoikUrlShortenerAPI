package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/big"
	"time"

	"github.com/AlphaExplorer0/StoikUrlShortenerAPI/repository"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

const shortLinkSize int = 7

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

type urlStorage interface {
	PutNewURL(ctx context.Context, originalUrl, shortUrl string) (string, error)
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

	res, err := ss.db.PutNewURL(ctx, url, tiny)

	// md5 can create collisions. If it's the case, we re-hash with an additional random payload until
	// url generated is unique. This can be improved with other heneration strategy like incremental counter.
	for err == repository.ErrURLAlreadyExists {
		ss.logger.Info(fmt.Sprintf("Rehashing for url %s", url))
		tiny = reHash(url, shortLinkSize)
		res, err = ss.db.PutNewURL(ctx, url, tiny)
	}

	if err != nil {
		return "", err
	}

	return "http://" + res, nil
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

func reHash(str string, size int) string {
	str = str + generateRandomString(10)
	return hashAndCut(str, size)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
