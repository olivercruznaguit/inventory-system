package service

import (
	"context"
	"errors"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type fakeUserAuthenticator struct {
	user model.User
	err  error
}

type fakeTokenGenerator struct {
	token string
	err   error
}

func (f *fakeUserAuthenticator) Authenticate(ctx context.Context, email, password string) (model.User, error) {
	if f.err != nil {
		return model.User{}, f.err
	}

	return f.user, nil
}

func (f *fakeTokenGenerator) GenerateToken(user model.User) (string, error) {
	if f.err != nil {
		return "", f.err
	}

	return f.token, nil
}

func TestAuthService_Authentication_Success(t *testing.T) {
	// ARRANGE
	expectedUser := model.User{
		ID:    1,
		Email: "test@gmail.com",
	}

	expectedToken := "test-jwt-token"

	userAuthenticator := &fakeUserAuthenticator{
		user: expectedUser,
	}

	tokenGenerator := &fakeTokenGenerator{
		token: expectedToken,
	}

	authService := NewAuthService(
		userAuthenticator,
		tokenGenerator,
	)

	// ACT
	token, err := authService.Login(
		context.Background(),
		expectedUser.Email,
		"testpassword",
	)

	// ASSERT
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != expectedToken {
		t.Errorf("expected token %q, got %q", expectedToken, token)
	}
}
func TestAuthService_Authentication_Fail(t *testing.T) {
	// ARRANGE
	expectedErr := errors.New("authentication failed")

	userAuthenticator := &fakeUserAuthenticator{
		err: expectedErr,
	}

	tokenGenerator := &fakeTokenGenerator{
		token: "should-not-be-generated",
	}

	authService := NewAuthService(userAuthenticator, tokenGenerator)

	// ACT
	token, err := authService.Login(
		context.Background(),
		"test@gmail.com",
		"wrongpassword",
	)

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if token != "" {
		t.Errorf("expected empty token, got %q", token)
	}
}
func TestAuthService_Token_Generator_Fail(t *testing.T) {
	// ARRANGE
	expectedErr := errors.New("token generation failed")

	user := model.User{
		ID:    1,
		Email: "test@gmail.com",
	}

	userAuthenticator := &fakeUserAuthenticator{
		user: user,
	}

	tokenGenerator := &fakeTokenGenerator{
		err: expectedErr,
	}

	authService := NewAuthService(
		userAuthenticator,
		tokenGenerator,
	)

	// ACT
	token, err := authService.Login(
		context.Background(),
		user.Email,
		"testpassword",
	)

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if token != "" {
		t.Errorf("expected empty token, got %q", token)
	}
}
