package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"web-scraper/internal/core/models"
	"web-scraper/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("✅ Initializing Seeding ...")
	godotenv.Load()

	db.ConnectDB()
	conn := db.GetDB()
	ctx := context.Background()

	users := []models.User{
		{Email: "alice@example.com", Name: "Alice Smith", Password: "Blue123!"},
		{Email: "bob@example.com", Name: "Bob Kent", Password: "Red543!$"},
		{Email: "alex@example2.com", Name: "Alex Lewindowski", Password: "Foo123"},
		{Email: "steve@example2.com", Name: "Steve Heyrmen", Password: "Foo123"},
		{Email: "carol@example2.com", Name: "Carol Askren", Password: "Foo123"},
	}

	for _, user := range users {
		var existingID string

		fmt.Println("Adding user: ", user.Email)
		fmt.Println("Checking if user exists... ")
		checkQuery := `SELECT id FROM users WHERE email = $1 LIMIT 1`
		err := conn.QueryRow(ctx, checkQuery, user.Email).Scan(&existingID)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("✅ User does not exist, adding now:", user.Email)

			hashed, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if hashErr != nil {
				log.Fatalf("❌ Failed to hash password for %s: %v", user.Email, hashErr)
			}

			user.Password = string(hashed)

			insertQuery := `
				INSERT INTO users (email, name, password)
				VALUES ($1, $2, $3)
				RETURNING id
			`

			var insertedID string
			err := conn.QueryRow(ctx, insertQuery,
				user.Email, user.Name, user.Password,
			).Scan(&insertedID)

			fmt.Println("inserted id: ", insertedID)

			if err != nil {
				log.Fatalf("❌ Failed to insert user %s: %v", user.Email, err)
			}

			log.Println("✅ Successfully added user:", user.Email)

		} else if err != nil {
			log.Fatalf("❌ Unexpected error during lookup: %v", err)
		} else {
			log.Println("⚠️  User already exists:", user.Email)
		}
	}

	fmt.Println("✅ Finished seeding data")
}
