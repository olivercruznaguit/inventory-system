package request

type CreateProductRequest struct {
	Name         string  `json:"name" binding:"required"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	CategoryID   *uint   `json:"categoryId"`
	MinimumStock int     `json:"minimumStock"`
}

type UpdateProductRequest struct {
	Name         string  `json:"name" binding:"required"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	Status       string  `json:"status" binding:"required"`
	CategoryID   *uint   `json:"categoryId"`
	MinimumStock int     `json:"minimumStock"`
}

type UpdateProductStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
