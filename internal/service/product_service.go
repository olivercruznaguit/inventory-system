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

func (ps *ProductService) GetProducts(ctx context.Context) ([]model.Product, error) {
	return ps.repository.GetProducts(ctx)
}

func (ps *ProductService) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	return ps.repository.GetProductByID(ctx, id)
}
