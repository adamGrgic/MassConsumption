package handlers

import (
	layout_vc "web-scraper/internal/components/layout"
	"web-scraper/internal/core/services"

	"github.com/gin-gonic/gin"
)

type ContentHandlers struct {
	Service services.ContentService
}

func NewContentHandlers() *ContentHandlers {
	return &ContentHandlers{}
}

func (h *ContentHandlers) GetLayout(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")

	model := layout_vc.Model{
		Context: c,
		Title:   "Home",
	}

	layout_vc.HTML(model).Render(c.Request.Context(), c.Writer)
}
