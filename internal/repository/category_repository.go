package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	var createdCategory model.Category
	var pgErr *pgconn.PgError

	err := cr.db.QueryRow(ctx,
		`INSERT INTO categories (name)
	 VALUES ($1) RETURNING
	 id,
	 name,
	 created_at,
	 updated_at`, category.Name).Scan(
		&createdCategory.ID,
		&createdCategory.Name,
		&createdCategory.CreatedAt,
		&createdCategory.UpdatedAt,
	)

	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Category{}, ErrCategoryAlreadyExists
		}

		return model.Category{}, fmt.Errorf("create category: %w", err)
	}

	return createdCategory, nil
}

func (cr *CategoryRepository) GetCategories(ctx context.Context) ([]model.Category, error) {
	rows, err := cr.db.Query(ctx,
		`SELECT 
		id, 
		name, 
		created_at, 
		updated_at
		FROM categories
		ORDER BY name`)

	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}

	defer rows.Close()

	var categories []model.Category

	for rows.Next() {
		var category model.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("get categories: %w", err)
		}

		categories = append(
			categories,
			category,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}

	return categories, nil
}

func (cr *CategoryRepository) GetCategoryByID(ctx context.Context, id int) (model.Category, error) {
	var category model.Category
	err := cr.db.QueryRow(ctx,
		`SELECT
		id,
		name,
		created_at,
		updated_at
		FROM categories WHERE id = $1`, id).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Category{}, ErrCategoryNotFound
		}

		return model.Category{}, fmt.Errorf("get category: %w", err)
	}

	return category, nil
}
