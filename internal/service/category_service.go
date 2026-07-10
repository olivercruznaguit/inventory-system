package service

import (
	"context"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category model.Category) (model.Category, error)
}

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repository CategoryRepository) *CategoryService {
	return &CategoryService{
		repository: repository,
	}
}

func (cr *CategoryService) CreateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	return cr.repository.CreateCategory(ctx, category)
}
