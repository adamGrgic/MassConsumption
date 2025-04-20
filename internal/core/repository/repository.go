package repository

import "github.com/jackc/pgx/v5/pgxpool"

func RegisterRepositories(db *pgxpool.Pool) {
	NewTaskRepository(db)
}
