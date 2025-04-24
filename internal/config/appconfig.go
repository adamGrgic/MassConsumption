package config

import (
	"log"
	"os"
	"web-scraper/internal/core/handlers"

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

	contentHandlers := handlers.NewContentHandlers()

	r.GET("/", contentHandlers.GetDashboardPage)

	return &App{
		Router: r,
	}
}
