package service

import (
	"context"
	"fmt"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category model.Category) (model.Category, error)
	GetCategories(ctx context.Context) ([]model.Category, error)
	GetCategoryByID(ctx context.Context, id int) (model.Category, error)
	UpdateCategory(ctx context.Context, category model.Category) (model.Category, error)
	DeleteCategory(ctx context.Context, id int) error
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

func (cs *CategoryService) GetCategories(ctx context.Context) ([]model.Category, error) {
	return cs.repository.GetCategories(ctx)
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
