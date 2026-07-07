package service

import (
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

func (ps *ProductService) GetProducts() []model.Product {
	products := ps.repository.GetProducts()
	activeProducts := make([]model.Product, 0, len(products))
	for _, product := range products {
		if product.Status == model.ProductStatusActive {
			activeProducts = append(activeProducts, product)
		}
	}
	return activeProducts
}
