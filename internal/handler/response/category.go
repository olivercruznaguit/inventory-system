package response

import (
	"time"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProductCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewCategoryResponse(category model.Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}
