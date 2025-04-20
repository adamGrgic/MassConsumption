package services

import (
	"context"
	"fmt"
	"goth-todo/internal/core/models"
	"goth-todo/internal/core/repository"
)

type UserService interface {
	GetUser(ctx context.Context, user *models.User, email string, password string) error
}

type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		Repo: userRepository,
	}
}

func (r *UserServiceImpl) GetUser(ctx context.Context, user *models.User, email string, password string) error {
	// Look up user

	if err := r.Repo.GetUser(ctx, user, email, password); err != nil {
		fmt.Println("Repo GetUser returned error: ", err)
		return err
	}

	return nil
}

func (r *UserServiceImpl) RegisterUser(email string, password string) {

}
