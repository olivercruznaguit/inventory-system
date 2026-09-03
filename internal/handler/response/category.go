package response

import (
	"time"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type CategoryResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	ProductCount int       `json:"productCount"`
}

type ProductCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CategoryListResponse struct {
	Data       []CategoryResponse `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

func NewCategoryResponse(category model.Category) CategoryResponse {
	return CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
		ProductCount: category.ProductCount,
	}
}
