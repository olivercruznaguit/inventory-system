package response

import (
	"time"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

func NewProductResponse(product model.Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Status:    string(product.Status),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
