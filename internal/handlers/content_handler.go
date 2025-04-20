package handlers

import (
	dashboard_vc "goth-todo/internal/components/dashboard"
	layout_vc "goth-todo/internal/components/layout"
	"goth-todo/internal/core/services"

	"github.com/gin-gonic/gin"
)

type ContentHandlers struct {
	Service services.ContentService
}

func NewContentHandlers() *ContentHandlers {
	return &ContentHandlers{}
}

// GetHomePage godoc
// @Summary      Get Home page
// @Description  Example of getting a parent component and a child component
// @Produce      html
// @Success      200 {string} string "HTML content"
// @Router       / [get]
func (h *ContentHandlers) GetDashboardPage(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")

	dashboard_model := dashboard_vc.Model{Context: c}

	model := layout_vc.Model{
		Context:   c,
		Title:     "Home",
		Component: dashboard_vc.HTML(dashboard_model),
	}

	layout_vc.HTML(model).Render(c.Request.Context(), c.Writer)
}
