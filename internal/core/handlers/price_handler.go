package handlers

import (
	"fmt"
	pricing_table_vc "web-scraper/internal/components/pricing-table"
	"web-scraper/internal/core/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type PriceHandler struct {
	PriceService services.PriceService
}

func NewPriceHandler(priceService services.PriceService) *PriceHandler {
	return &PriceHandler{
		PriceService: priceService,
	}
}

func (p *PriceHandler) GetPrices(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")

	prices, err := p.PriceService.GetPrices()
	if err != nil {
		fmt.Println("Could not retrieve prices: ", err)
		return
	}

	log.Info().
		Interface("prices", prices).
		Msg("Successfully retrieved prices")

	model := pricing_table_vc.Model{
		PriceRecords: prices,
	}

	pricing_table_vc.HTML(model).Render(c.Request.Context(), c.Writer)
}
