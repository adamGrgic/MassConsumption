package main

import (
	"flag"
	"fmt"
	"goth-todo/internal/db"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool
var dir = "./internal/db/migrations"

// use this instead of makefile for more complex workflows
func main() {

	MIGRATE_DB := flag.Bool("migrate-db", false, "apply migration from app startup command")
	MIGRATE_CLEAN := flag.Bool("migrate-clean", false, "clean migration directory")

	flag.Parse()

	switch {
	case *MIGRATE_DB:
		db.ConnectDB()
	case *MIGRATE_CLEAN:
		var input string
		fmt.Println("Cleaning your migration directory will destroy all of your migrations. Are you sure? (y/N)")
		fmt.Scanln(&input)

		if input == "y" || input == "Y" {
			fmt.Println("✅ Confirmed!")
			entries, err := os.ReadDir(dir)
			if err != nil {
				log.Panic("❌ could not read migration directory: ", dir)
			}

			for _, entry := range entries {
				entryPath := filepath.Join(dir, entry.Name())
				err := os.RemoveAll(entryPath)
				if err != nil {
					log.Panic("❌ could not read migration entry: ", entryPath)
				}
			}

			fmt.Println("✅ Migrations cleared!")
		} else {
			fmt.Println("❌ Aborted.")
		}
	}

}
