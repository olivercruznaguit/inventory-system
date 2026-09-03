package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func (cr *CategoryRepository) GetCategories(ctx context.Context, filter model.CategoryFilter) ([]model.Category, error) {
	pagination := filter.Pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	limit := pagination.PageSize

	var queryParts []string
	var conditions []string
	var args []any

	queryParts = append(queryParts,
		`SELECT 
		c.id, 
		c.name, 
		c.created_at, 
		c.updated_at,
		COUNT(p.id) AS product_count 
		FROM categories c LEFT JOIN products p ON c.id = p.category_id`)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")

		conditions = append(
			conditions,
			fmt.Sprintf("c.name ILIKE $%d", len(args)),
		)
	}

	if len(conditions) > 0 {
		queryParts = append(
			queryParts,
			"WHERE "+strings.Join(conditions, " AND "),
		)
	}

	queryParts = append(queryParts, "GROUP BY c.id, c.name, c.created_at, c.updated_at")
	queryParts = append(queryParts, fmt.Sprintf("ORDER BY %s %s", filter.SortBy, filter.SortOrder))

	args = append(args, limit)
	queryParts = append(queryParts, fmt.Sprintf("LIMIT $%d", len(args)))

	args = append(args, offset)
	queryParts = append(queryParts, fmt.Sprintf("OFFSET $%d", len(args)))

	queryString := strings.Join(queryParts, " ")

	rows, err := cr.db.Query(ctx, queryString, args...)

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
			&category.ProductCount,
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

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	var updatedCategory model.Category
	err := cr.db.QueryRow(ctx,
		`UPDATE categories SET 
		 name = $1,
		 updated_at = NOW() 
		WHERE id = $2 
		RETURNING
		id,
		name,
		created_at,
		updated_at`, category.Name, category.ID).Scan(
		&updatedCategory.ID,
		&updatedCategory.Name,
		&updatedCategory.CreatedAt,
		&updatedCategory.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Category{}, ErrCategoryNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Category{}, ErrCategoryAlreadyExists
		}

		return model.Category{}, fmt.Errorf("update category: %w", err)
	}

	return updatedCategory, nil
}

func (cr *CategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	result, err := cr.db.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)

	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (cr *CategoryRepository) CountCategories(ctx context.Context, filter model.CategoryFilter) (int, error) {
	var count int
	var queryParts []string
	var conditions []string
	var args []any

	queryParts = append(queryParts, "SELECT COUNT(*) FROM categories")

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}

	if len(conditions) > 0 {
		queryParts = append(
			queryParts,
			"WHERE "+strings.Join(conditions, " AND "),
		)
	}

	queryString := strings.Join(queryParts, " ")

	err := cr.db.QueryRow(ctx, queryString, args...).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("count categories: %w", err)
	}

	return count, nil
}
