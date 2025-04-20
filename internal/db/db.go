package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var DB *pgxpool.Pool

// ConnectDB initializes the pgx pool connection
func ConnectDB() {
	dsn := os.Getenv("DB_URL") // or hardcode for now
	if dsn == "" {
		log.Fatal().Msg("No DB environment variable setup")
	}

	log.Info().Str("dsn", dsn).Msg("database url")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal().Msg("❌ Failed to connect to the database")
	}

	if err := DB.Ping(ctx); err != nil {
		log.Fatal().Msg("❌ Unable to ping DB")
	}

	fmt.Println("✅ Connected to PostgreSQL with pgx")
}

func GetDB() *pgxpool.Pool {
	return DB
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Info().Msg("Database connection closed")
	}
}
