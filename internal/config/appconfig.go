package config

import (
	"log"
	"os"
	"web-scraper/internal/core/handlers"
	"web-scraper/internal/core/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	Router *gin.Engine
}

func loadEnv() {
	// Optional: only load .env in local environments
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, assuming CI/production")
		}
	}
}

func RunApp() *App {

	loadEnv()

	r := gin.Default()
	r.Static("/public", "./public")

	// Initialize repositories

	priceService := services.NewPriceService()

	contentHandlers := handlers.NewContentHandlers()
	priceHandlers := handlers.NewPriceHandler(priceService)

	r.GET("/", contentHandlers.GetLayout)

	r.GET("/prices-table/get", priceHandlers.GetPrices)

	return &App{
		Router: r,
	}
}
