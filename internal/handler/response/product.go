package response

import (
	"time"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type ProductResponse struct {
	ID           uint                     `json:"id"`
	Name         string                   `json:"name"`
	Price        float64                  `json:"price"`
	Status       string                   `json:"status"`
	Quantity     int                      `json:"quantity"`
	MinimumStock int                      `json:"minimumStock"`
	Category     *ProductCategoryResponse `json:"category,omitempty"`
	CreatedAt    time.Time                `json:"createdAt"`
	UpdatedAt    time.Time                `json:"updatedAt"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type ProductListResponse struct {
	Data       []ProductResponse  `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type StockResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	MinimumStock int    `json:"minimumStock"`
}

func NewProductResponse(product model.Product) ProductResponse {
	productResponse := ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Price:        product.Price,
		Status:       string(product.Status),
		Quantity:     product.Quantity,
		MinimumStock: product.MinimumStock,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	if product.Category != nil {
		productResponse.Category = &ProductCategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		}
	}

	return productResponse
}

func NewStockResponse(product model.Product) StockResponse {
	return StockResponse{
		ID:           product.ID,
		Name:         product.Name,
		Quantity:     product.Quantity,
		MinimumStock: product.MinimumStock,
	}
}
