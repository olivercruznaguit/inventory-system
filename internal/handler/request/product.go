package request

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name   string  `json:"name" binding:"required"`
	Price  float64 `json:"price" binding:"required,gt=0"`
	Status string  `json:"status" binding:"required"`
}

type UpdateProductStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
