package service

import (
	"context"
	"fmt"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (us *UserService) Authenticate(ctx context.Context, email string, password string) (model.User, error) {
	if !isValidEmail(email) {
		return model.User{}, fmt.Errorf("authenticate user: %w", ErrInvalidEmailAddress)
	}

	user, err := us.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return model.User{}, fmt.Errorf("authenticate user: %w", ErrIncorrectPassword)
	}

	return user, nil
}

func (us *UserService) Register(ctx context.Context, email string, password string) (model.User, error) {
	if !isValidEmail(email) {
		return model.User{}, fmt.Errorf("register user: %w", ErrInvalidEmailAddress)
	}

	// hash pass
	passwordHash, err := HashPassword(password)
	if err != nil {
		return model.User{}, fmt.Errorf("register user: %w", err)
	}

	user, err := us.repository.CreateUser(
		ctx,
		model.User{
			Email:        email,
			PasswordHash: passwordHash,
		},
	)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
