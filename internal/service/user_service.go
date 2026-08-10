package service

import (
	"context"
	"fmt"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (us *UserService) VerifyUser(ctx context.Context, email string, password string) (bool, error) {
	if !isValidEmail(email) {
		return false, fmt.Errorf("verify user: %w", ErrInvalidEmailAddress)
	}

	user, err := us.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return false, fmt.Errorf("verify user: %w", ErrIncorrectPassword)
	}

	return true, nil
}
