package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category model.Category) (model.Category, error)
	GetCategories(ctx context.Context, filter model.CategoryFilter) ([]model.Category, error)
	GetCategoryByID(ctx context.Context, id int) (model.Category, error)
	UpdateCategory(ctx context.Context, category model.Category) (model.Category, error)
	DeleteCategory(ctx context.Context, id int) error
	CountCategories(ctx context.Context, filter model.CategoryFilter) (int, error)
}

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repository CategoryRepository) *CategoryService {
	return &CategoryService{
		repository: repository,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	return cs.repository.CreateCategory(ctx, category)
}

func (cs *CategoryService) GetCategories(ctx context.Context, filter model.CategoryFilter) (model.CategoryList, error) {
	if filter.Pagination.Page < 1 {
		filter.Pagination.Page = 1
	}

	if filter.Pagination.PageSize < 1 {
		filter.Pagination.PageSize = 20
	}

	if filter.Pagination.PageSize > 100 {
		filter.Pagination.PageSize = 100
	}

	if filter.SortBy == "" {
		filter.SortBy = "id"
	}

	filter.SortBy = strings.ToLower(filter.SortBy)
	switch filter.SortBy {
	case "id", "name":
		//valid
	default:
		return model.CategoryList{},
			fmt.Errorf("get categories: %w", ErrInvalidSortBy)
	}

	if filter.SortOrder == "" {
		filter.SortOrder = "asc"
	}

	switch strings.ToLower(filter.SortOrder) {
	case "asc", "desc":
		filter.SortOrder = strings.ToUpper(filter.SortOrder)
	default:
		return model.CategoryList{},
			fmt.Errorf("get categories: %w", ErrInvalidSortOrder)
	}

	categories, err := cs.repository.GetCategories(ctx, filter)
	if err != nil {
		return model.CategoryList{}, fmt.Errorf("get categories: %w", err)
	}

	totalItems, err := cs.repository.CountCategories(ctx, filter)
	if err != nil {
		return model.CategoryList{}, fmt.Errorf("count categories: %w", err)
	}

	totalPages := (totalItems + filter.Pagination.PageSize - 1) / filter.Pagination.PageSize

	filter.Pagination.TotalItems = totalItems
	filter.Pagination.TotalPages = totalPages

	return model.CategoryList{
		Categories: categories,
		Pagination: filter.Pagination,
	}, nil
}

func (cs *CategoryService) GetCategoryByID(ctx context.Context, id int) (model.Category, error) {
	return cs.repository.GetCategoryByID(ctx, id)
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, category model.Category) (model.Category, error) {

	updatedCategory, err := cs.repository.UpdateCategory(ctx, category)

	if err != nil {
		return model.Category{}, fmt.Errorf("update category: %w", err)
	}

	return updatedCategory, nil
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return cs.repository.DeleteCategory(ctx, id)
}
