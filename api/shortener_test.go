package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_Should_Return_HTTP200_When_payload_Is_Valid(t *testing.T) {
	sh := ShortenerHandler{Logger: zap.NewNop()}
	engine := gin.New()
	engine.POST("/", sh.Handle)

	response := do(engine, "/", http.MethodPost, `{
		"long_url": "http://test"
	}`)

	assert.Equal(t, http.StatusOK, response.Code)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"short_url":"http://ok"
	}`, string(body))
}

func Test_Should_Return_HTTP400_When_payload_Is_Invalid(t *testing.T) {
	sh := ShortenerHandler{Logger: zap.NewNop()}
	engine := gin.New()
	engine.POST("/", sh.Handle)

	response := do(engine, "/", http.MethodPost, `{
		"toto": "something"
	}`)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func do(app http.Handler, url string, method string, body string) *httptest.ResponseRecorder {
	bodyReader := strings.NewReader(body)
	request := httptest.NewRequest(method, url, bodyReader)
	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)
	return recorder
}