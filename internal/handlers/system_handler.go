package handlers

import (
	"goth-todo/internal/core/services"
	"net/http"

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
