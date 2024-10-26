package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func Test_Should_Return_HTTP301_When_payload_Is_Valid(t *testing.T) {

	mockedRedirector := new(mockRedirectService)
	mockedRedirector.On("FindURL", mock.Anything, mock.Anything).Return("https://fr.wikipedia.org", nil)

	sh := RedirectHandler{Logger: zap.NewNop(), Service: mockedRedirector}

	engine := gin.New()
	engine.GET("/short.io/:key", sh.Handle)

	response := do(engine, "/short.io/abcdefg", http.MethodGet, "")

	assert.Equal(t, http.StatusMovedPermanently, response.Code)
	mockedRedirector.AssertExpectations(t)
}

func Test_Should_Return_HTTP500_When_url_not_found(t *testing.T) {
	mockedRedirector := new(mockRedirectService)
	mockedRedirector.On("FindURL", mock.Anything, mock.Anything).Return("", errors.New("oops"))

	sh := RedirectHandler{Logger: zap.NewNop(), Service: mockedRedirector}

	engine := gin.New()
	engine.GET("/short.io/:key", sh.Handle)

	response := do(engine, "/short.io/abcdefg", http.MethodGet, "")

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	mockedRedirector.AssertExpectations(t)
}

type mockRedirectService struct {
	mock.Mock
}

func (mock *mockRedirectService) FindURL(ctx context.Context, shortKey string) (string, error) {
	args := mock.Called(ctx, shortKey)
	return args.Get(0).(string), args.Error(1)
}
