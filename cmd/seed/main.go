package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"goth-todo/internal/core/models"
	"goth-todo/internal/db"

	"github.com/google/uuid"
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

	accountIDs := make(map[string]uuid.UUID)

	accounts := []models.Account{
		{Name: "Test Account 1"},
		{Name: "Test Account 2"},
	}

	users := []models.User{
		{AccountId: accountIDs["Test Account 1"], Email: "alice@example.com", FirstName: "Alice", LastName: "Smith", Password: "Blue123!"},
		{AccountId: accountIDs["Test Account 1"], Email: "bob@example.com", FirstName: "Bob", LastName: "Kent", Password: "Red543!$"},
		{AccountId: accountIDs["Test Account 2"], Email: "alex@example2.com", FirstName: "Alex", LastName: "Lewindowksi", Password: "Foo123"},
		{AccountId: accountIDs["Test Account 2"], Email: "steve@example2.com", FirstName: "Steve", LastName: "Heyrman", Password: "Foo123"},
		{AccountId: accountIDs["Test Account 2"], Email: "carol@example2.com", FirstName: "Carol", LastName: "Askren", Password: "Foo123"},
	}

	for _, account := range accounts {
		var existingID string

		checkQuery := `SELECT id FROM accounts WHERE name = $1 LIMIT 1`
		err := conn.QueryRow(ctx, checkQuery, account.Name).Scan(&existingID)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("✅ Account does not exist, adding now:", account.Name)

			var insertedID string // <- define this before using it
			insertQuery := `
				INSERT INTO accounts (name)
				VALUES ($1)
				RETURNING id
			`
			err := conn.QueryRow(ctx, insertQuery, account.Name).Scan(&insertedID)
			if err != nil {
				log.Fatalf("❌ Failed to insert account %s: %v", account.Name, err)
			}

			log.Printf("✅ Successfully added account: %s (id: %s)", account.Name, insertedID)

			accountIDs[account.Name] = account.ID

		} else if err != nil {
			log.Fatalf("❌ Unexpected error during lookup for account %s: %v", account.Name, err)
		} else {
			log.Printf("⚠️  Account already exists: %s (id: %s)", account.Name, existingID)
		}
	}

	for _, user := range users {
		var existingID string

		fmt.Println("Adding user: ", user.Email)
		fmt.Println("Cehcking if user exists... ")
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
				INSERT INTO users (email, first_name, last_name, password, account_id)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`

			var insertedID string
			err := conn.QueryRow(ctx, insertQuery,
				user.Email, user.FirstName, user.LastName, user.Password, user.AccountId,
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

func generateUUID() uuid.UUID {
	panic("unimplemented")
}
