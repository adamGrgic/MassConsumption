package main

import (
	"os"
	"web-scraper/internal/config"
	"web-scraper/internal/logging"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logging.ConfigureLogging()

	app := config.RunApp()

	app.Router.Run(os.Getenv("PORT"))
}
