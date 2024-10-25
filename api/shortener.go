package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShortenRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
}

type ShortenerService interface {
	ShortenURL(ctx context.Context, url string) (string, error)
}

type ShortenerHandler struct {
	Logger  *zap.Logger
	Service ShortenerService
}

func (sh *ShortenerHandler) Handle(c *gin.Context) {
	var body ShortenRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		sh.Logger.Error(fmt.Sprintf("could not unbind request body: %s", err.Error()))
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_, err := url.Parse(body.LongUrl)
	if err != nil {
		sh.Logger.Error(fmt.Sprintf("could not parse input url: %s", err.Error()))
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	tiny, err := sh.Service.ShortenURL(c.Request.Context(), body.LongUrl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.JSON(http.StatusOK, map[string]string{"short_url": tiny})
}
