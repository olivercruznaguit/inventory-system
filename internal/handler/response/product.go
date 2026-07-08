package response

type ProductResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
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
