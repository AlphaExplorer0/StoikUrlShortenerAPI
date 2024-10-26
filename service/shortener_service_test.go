package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func Test_hashAndCut_produces_expected_output_format(t *testing.T) {
	url := "https://fr.wikipedia.org/wiki/Uniform_Resource_Locator"
	tiny := hashAndCut(url, 7)
	fmt.Printf("%s\n", tiny)
	assert.Equal(t, 7, len(tiny))
}

func Test_ShortenURL_returns_expected_output(t *testing.T) {
	mockedDb := new(mockUrlStorage)
	mockedDb.On("PutNewURL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ss := shortenerService{logger: zap.NewNop(), db: mockedDb}

	url, err := ss.ShortenURL(context.Background(), "https://fr.wikipedia.org")

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/short.io/", url[:len(url)-shortLinkSize])
	mockedDb.AssertExpectations(t)
}

func Test_ShortenURL_returns_error_if_dbInsertion_fails(t *testing.T) {
	mockedDb := new(mockUrlStorage)
	mockedDb.On("PutNewURL", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("fail"))

	ss := shortenerService{logger: zap.NewNop(), db: mockedDb}

	_, err := ss.ShortenURL(context.Background(), "https://fr.wikipedia.org")

	assert.Error(t, err)
	assert.Equal(t, "fail", err.Error())
	mockedDb.AssertExpectations(t)
}

type mockUrlStorage struct {
	mock.Mock
}

func (mock *mockUrlStorage) PutNewURL(ctx context.Context, originalUrl, shortUrl string) error {
	args := mock.Called(ctx, originalUrl, shortUrl)
	return args.Error(0)
}
