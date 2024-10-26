package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func Test_FindURL_returns_expected_output(t *testing.T) {
	mockedDb := new(mockUrlGetter)
	mockedDb.On("GetOriginalURL", mock.Anything, mock.Anything).Return("https://fr.wikipedia.org", nil)

	ss := redirectService{logger: zap.NewNop(), db: mockedDb}

	url, err := ss.FindURL(context.Background(), "abcdefg")

	assert.NoError(t, err)
	assert.Equal(t, "https://fr.wikipedia.org", url)
	mockedDb.AssertExpectations(t)
}

func Test_FindURL_returns_error_if_dbInsertion_fails(t *testing.T) {
	mockedDb := new(mockUrlGetter)
	mockedDb.On("GetOriginalURL", mock.Anything, mock.Anything).Return("", errors.New("fail"))

	ss := redirectService{logger: zap.NewNop(), db: mockedDb}

	_, err := ss.FindURL(context.Background(), "abcdefg")

	assert.Error(t, err)
	assert.Equal(t, "fail", err.Error())
	mockedDb.AssertExpectations(t)
}

type mockUrlGetter struct {
	mock.Mock
}

func (mock *mockUrlGetter) GetOriginalURL(ctx context.Context, key string) (string, error) {
	args := mock.Called(ctx, key)
	return args.Get(0).(string), args.Error(1)
}
