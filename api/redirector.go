package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RedirectService interface {
	FindURL(ctx context.Context, shortKey string) (string, error)
}

type RedirectHandler struct {
	Logger  *zap.Logger
	Service RedirectService
}

func (rh *RedirectHandler) Handle(c *gin.Context) {
	shortKey := c.Param("key")

	if shortKey == "" {
		rh.Logger.Error("shortened key is missing")
		c.JSON(http.StatusBadRequest, map[string]string{"error": "shortened key is missing"})
		return
	}

	originalURL, err := rh.Service.FindURL(c.Request.Context(), shortKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
