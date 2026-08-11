package service

import (
	"context"
	"fmt"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type AuthService struct {
	userService  UserAuthenticator
	tokenService TokenGenerator
}

type UserAuthenticator interface {
	Authenticate(ctx context.Context, email, password string) (model.User, error)
}

type TokenGenerator interface {
	GenerateToken(user model.User) (string, error)
}

func NewAuthService(userService UserAuthenticator, tokenService TokenGenerator) *AuthService {
	return &AuthService{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userService.Authenticate(ctx, email, password)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	token, err := s.tokenService.GenerateToken(user)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	return token, nil
}
