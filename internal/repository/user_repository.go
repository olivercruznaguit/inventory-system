package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	row := ur.db.QueryRow(ctx, `
	SELECT 
		id,
		email,
		password_hash,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	LIMIT 1
	`, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}
