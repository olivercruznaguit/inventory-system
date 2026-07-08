package service

import (
	"context"

	"github.com/olivercruznaguit/inventory-system/internal/model"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
)

type ProductService struct {
	repository *repository.ProductRepository
}

func NewProductService(repository *repository.ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (ps *ProductService) GetProducts(ctx context.Context, pagination model.Pagination) ([]model.Product, error) {
	if pagination.Page < 1 {
		pagination.Page = 1
	}

	if pagination.PageSize < 1 {
		pagination.PageSize = 20
	}

	if pagination.PageSize > 100 {
		pagination.PageSize = 100
	}

	return ps.repository.GetProducts(ctx, pagination)
}

func (ps *ProductService) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	return ps.repository.GetProductByID(ctx, id)
}

func (ps *ProductService) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	return ps.repository.CreateProduct(ctx, product)
}

func (ps *ProductService) UpdateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	return ps.repository.UpdateProduct(ctx, product)
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return ps.repository.DeleteProduct(ctx, id)
}
