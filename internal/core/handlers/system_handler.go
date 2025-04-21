package handlers

import (
	"net/http"
	"web-scraper/internal/core/services"

	"github.com/gin-gonic/gin"
)

type SystemHandlers struct {
	Service services.SystemService
}

func NewSystemHandlers(service services.SystemService) *SystemHandlers {
	return &SystemHandlers{Service: service}
}

func (h *SystemHandlers) Ping(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "PONG",
	})
}
