package request

type StockRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}
