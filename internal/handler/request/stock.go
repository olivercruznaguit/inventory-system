package request

type StockRequest struct {
	Quantity int     `json:"quantity" binding:"required"`
	Reason   *string `json:"reason"`
}
