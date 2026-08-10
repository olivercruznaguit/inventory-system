package service

import (
	"context"
	"errors"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/model"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
)

type fakeUserRepository struct {
	getUserByEmailCalls int
	err                 error
	user                model.User
}

func (f *fakeUserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	f.getUserByEmailCalls++

	if f.err != nil {
		return model.User{}, f.err
	}

	return f.user, nil
}

func TestUserService_VerifyUser_Success(t *testing.T) {
	// ARRANGE
	password := "testpassword"
	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repository := &fakeUserRepository{
		user: model.User{
			Email:        "test@gmail.com",
			PasswordHash: passwordHash,
		},
	}

	service := NewUserService(repository)

	// ACT
	result, err := service.VerifyUser(context.Background(), "test@gmail.com", password)

	// ASSERT
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !result {
		t.Errorf("expected result %t, got %v", true, result)
	}

	if repository.getUserByEmailCalls != 1 {
		t.Errorf(
			"expected repository to be called once, got %d calls",
			repository.getUserByEmailCalls,
		)
	}
}

func TestUserService_VerifyUser_InvalidEmail(t *testing.T) {
	// ARRANGE
	repository := &fakeUserRepository{}

	service := NewUserService(repository)

	// ACT
	result, err := service.VerifyUser(context.Background(), "test@gmail.", "testpassword")

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrInvalidEmailAddress) {
		t.Errorf("expected error %s, got %s", ErrInvalidEmailAddress, err.Error())
	}

	if result {
		t.Errorf("expected result %t, got %v", false, result)
	}

	if repository.getUserByEmailCalls != 0 {
		t.Errorf(
			"expected repository not to be called, got %d calls",
			repository.getUserByEmailCalls,
		)
	}
}

func TestUserService_VerifyUser_IncorrectPassword(t *testing.T) {
	// ARRANGE
	password := "testpassword"
	email := "test@gmail.com"
	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repository := &fakeUserRepository{
		user: model.User{
			Email:        email,
			PasswordHash: passwordHash,
		},
	}

	service := NewUserService(repository)

	// ACT
	result, err := service.VerifyUser(context.Background(), email, "passwordtest")

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrIncorrectPassword) {
		t.Errorf("expected error %s, got %s", ErrIncorrectPassword, err.Error())
	}

	if result {
		t.Errorf("expected result %t, got %v", false, result)
	}

	if repository.getUserByEmailCalls != 1 {
		t.Errorf(
			"expected repository to be called once, got %d calls",
			repository.getUserByEmailCalls,
		)
	}
}

func TestUserService_VerifyUser_UserNotFound(t *testing.T) {
	// ARRANGE
	userRepository := &fakeUserRepository{
		err: repository.ErrUserNotFound,
	}

	service := NewUserService(userRepository)

	// ACT
	result, err := service.VerifyUser(context.Background(), "test@gmail.com", "passwordtest")

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Errorf("expected error %s, got %s", repository.ErrUserNotFound, err.Error())
	}

	if result {
		t.Errorf("expected result %t, got %v", false, result)
	}

	if userRepository.getUserByEmailCalls != 1 {
		t.Errorf(
			"expected repository to be called once, got %d calls",
			userRepository.getUserByEmailCalls,
		)
	}
}
