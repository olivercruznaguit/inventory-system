package model

type StockRequest struct {
	ProductID int
	Quantity  int
	Reason    *string
}
