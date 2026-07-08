package service

import (
	"context"
	"fmt"

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

func (ps *ProductService) GetProducts(ctx context.Context, pagination model.Pagination) (model.ProductList, error) {
	if pagination.Page < 1 {
		pagination.Page = 1
	}

	if pagination.PageSize < 1 {
		pagination.PageSize = 20
	}

	if pagination.PageSize > 100 {
		pagination.PageSize = 100
	}

	products, err := ps.repository.GetProducts(ctx, pagination)
	if err != nil {
		return model.ProductList{}, fmt.Errorf("get products: %w", err)
	}

	totalItems, err := ps.repository.CountProducts(ctx)
	if err != nil {
		return model.ProductList{}, fmt.Errorf("count products: %w", err)
	}

	totalPages := (totalItems + pagination.PageSize - 1) / pagination.PageSize

	pagination.TotalItems = totalItems
	pagination.TotalPages = totalPages

	return model.ProductList{
		Products:   products,
		Pagination: pagination,
	}, nil
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
