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
	createUserCalls     int
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

func (f *fakeUserRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	f.getUserByEmailCalls++

	if f.err != nil {
		return model.User{}, f.err
	}

	return f.user, nil
}

func TestUserService_Authenticate_Success(t *testing.T) {
	// ARRANGE
	password := "testpassword"
	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repository := &fakeUserRepository{
		user: model.User{
			ID:           1,
			Email:        "test@gmail.com",
			PasswordHash: passwordHash,
		},
	}

	service := NewUserService(repository)

	// ACT
	user, err := service.Authenticate(context.Background(), "test@gmail.com", password)

	// ASSERT
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != repository.user.ID {
		t.Errorf("expected user ID %d, got %d", repository.user.ID, user.ID)
	}

	if user.Email != repository.user.Email {
		t.Errorf("expected email %s, got %s", repository.user.Email, user.Email)
	}

	if repository.getUserByEmailCalls != 1 {
		t.Errorf(
			"expected repository to be called once, got %d calls",
			repository.getUserByEmailCalls,
		)
	}
}

func TestUserService_Authenticate_InvalidEmail(t *testing.T) {
	// ARRANGE
	repository := &fakeUserRepository{}

	service := NewUserService(repository)

	// ACT
	_, err := service.Authenticate(context.Background(), "test@gmail.", "testpassword")

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrInvalidEmailAddress) {
		t.Errorf("expected error %s, got %s", ErrInvalidEmailAddress, err.Error())
	}

	if repository.getUserByEmailCalls != 0 {
		t.Errorf(
			"expected repository not to be called, got %d calls",
			repository.getUserByEmailCalls,
		)
	}
}

func TestUserService_Authenticate_IncorrectPassword(t *testing.T) {
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
	_, err = service.Authenticate(context.Background(), email, "passwordtest")

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrIncorrectPassword) {
		t.Errorf("expected error %s, got %s", ErrIncorrectPassword, err.Error())
	}

	if repository.getUserByEmailCalls != 1 {
		t.Errorf(
			"expected repository to be called once, got %d calls",
			repository.getUserByEmailCalls,
		)
	}
}

func TestUserService_Authenticate_UserNotFound(t *testing.T) {
	// ARRANGE
	userRepository := &fakeUserRepository{
		err: repository.ErrUserNotFound,
	}

	service := NewUserService(userRepository)

	// ACT
	_, err := service.Authenticate(context.Background(), "test@gmail.com", "passwordtest")

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Errorf("expected error %s, got %s", repository.ErrUserNotFound, err.Error())
	}

	if userRepository.getUserByEmailCalls != 1 {
		t.Errorf(
			"expected repository to be called once, got %d calls",
			userRepository.getUserByEmailCalls,
		)
	}
}
