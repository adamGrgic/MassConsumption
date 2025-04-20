package repository

import (
	"context"
	"errors"
	"fmt"
	"goth-todo/internal/core/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUser(ctx context.Context, user *models.User, username string, password string) error {
	fmt.Println("UserRepository search user: ", username)

	query := `
		SELECT id, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	row := r.DB.QueryRow(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		fmt.Println("User not found or error: ", err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Invalid password")
		return errors.New("invalid credentials")
	}

	fmt.Println("User successfully validated: ", user.Email)
	return nil
}

// func (r *UserRepository) GetUser(user models.User, email string, password string) error {
// 	// Look up user

// 	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
// 		return err
// 	}

// 	// Check password (assuming it's hashed)
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
// 		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
// 		return err
// 	}
// }

// func (r *TaskRepository) GetTasks() ([]models.Task, error) {
// 	var tasks []models.Task
// 	err := r.DB.Find(&tasks).Error
// 	return tasks, err
// }

// func (r *TaskRepository) AddTask(task *models.Task) error {
// 	return r.DB.Create(task).Error
// }

// func (r *TaskRepository) ToggleTask(id string) error {
// 	var task models.Task
// 	err := r.DB.First(&task, "id = ?", id).Error
// 	if err != nil {
// 		return err
// 	}
// 	return r.DB.Save(&task).Error
// }

// func (r *TaskRepository) DeleteTask(id string) {
// 	r.DB.Delete(&models.Task{}, id)
// }
