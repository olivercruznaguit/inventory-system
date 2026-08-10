package inventory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/integration/testutil"
	"github.com/olivercruznaguit/inventory-system/internal/model"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
)

func TestUserRepository_GetUserByEmail_Success(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	user := testutil.SeedUser(t, ctx, db, model.User{
		Email:        "test@gmail.com",
		PasswordHash: "test-hash/123",
	})

	userRepository := repository.NewUserRepository(db.DB())

	// ACT
	fetchedUser, err := userRepository.GetUserByEmail(ctx, user.Email)

	// ASSERT
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != fetchedUser.Email {
		t.Fatalf("expected user email %s, got %s", user.Email, fetchedUser.Email)
	}
}

func TestUserRepository_GetUserByEmail_NotFound(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	testutil.SeedUser(t, ctx, db, model.User{
		Email:        "test@gmail.com",
		PasswordHash: "test-hash/123",
	})

	userRepository := repository.NewUserRepository(db.DB())

	// ACT
	_, err := userRepository.GetUserByEmail(ctx, "test@yahoo.com")

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected error %v, got %v", repository.ErrUserNotFound, err)
	}
}
