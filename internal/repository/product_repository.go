package repository

import (
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type ProductRepository struct{}

func (r ProductRepository) GetProducts() []model.Product {

	products := []model.Product{
		{
			ID:     1,
			Name:   "Laptop",
			Price:  10.0,
			Status: model.ProductStatusActive,
		},
		{
			ID:     2,
			Name:   "Keyboard",
			Price:  10.0,
			Status: model.ProductStatusActive,
		},
		{
			ID:     3,
			Name:   "Mouse",
			Price:  5.0,
			Status: model.ProductStatusInactive,
		}}

	return products
}
