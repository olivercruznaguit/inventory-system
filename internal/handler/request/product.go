package request

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	CategoryID *uint   `json:"categoryId"`
}

type UpdateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Status     string  `json:"status" binding:"required"`
	CategoryID *uint   `json:"categoryId"`
}

type UpdateProductStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
