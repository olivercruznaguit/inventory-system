package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olivercruznaguit/inventory-system/internal/handler/request"
	"github.com/olivercruznaguit/inventory-system/internal/handler/response"
	"github.com/olivercruznaguit/inventory-system/internal/model"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
	"github.com/olivercruznaguit/inventory-system/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	ctx := c.Request.Context()

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	filterStatusStr := c.Query("status")
	search := c.Query("search")
	categoryIdStr := c.Query("categoryId")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	var categoryID *uint
	if categoryIdStr != "" {
		categoryIdParsed, err := strconv.Atoi(categoryIdStr)
		if err != nil || categoryIdParsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		categoryIDCasted := uint(categoryIdParsed)
		categoryID = &categoryIDCasted
	}

	pagination := model.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	filter := model.ProductFilter{
		Status:     model.ProductStatus(filterStatusStr),
		Search:     search,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
		Pagination: pagination,
		CategoryID: categoryID,
	}

	products, err := h.service.GetProducts(ctx, filter)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProductStatus) ||
			errors.Is(err, service.ErrInvalidSortBy) ||
			errors.Is(err, service.ErrInvalidSortOrder) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productResponses []response.ProductResponse
	for _, product := range products.Products {
		productResponses = append(productResponses, response.NewProductResponse(product))
	}

	productListResponse := response.ProductListResponse{
		Data: productResponses,
		Pagination: response.PaginationResponse{
			Page:       products.Pagination.Page,
			PageSize:   products.Pagination.PageSize,
			TotalItems: products.Pagination.TotalItems,
			TotalPages: products.Pagination.TotalPages,
		},
	}

	c.JSON(http.StatusOK, productListResponse)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.service.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.NewProductResponse(product))
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	product := model.Product{
		Name:         req.Name,
		Price:        req.Price,
		MinimumStock: req.MinimumStock,
	}

	if req.CategoryID != nil {
		product.Category = &model.Category{
			ID: *req.CategoryID,
		}
	}

	createdProduct, err := h.service.CreateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProductMinimumStock) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product minimum stock"})
			return
		}

		if errors.Is(err, repository.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	product := model.Product{
		ID:           uint(id),
		Name:         req.Name,
		Price:        req.Price,
		MinimumStock: req.MinimumStock,
		Status:       model.ProductStatus(req.Status),
	}

	if req.CategoryID != nil {
		product.Category = &model.Category{
			ID: *req.CategoryID,
		}
	}

	updatedProduct, err := h.service.UpdateProduct(ctx, product)
	if err != nil {
		switch {

		case errors.Is(err, service.ErrInvalidProductMinimumStock):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid product minimum stock",
			})

		case errors.Is(err, service.ErrInvalidProductStatus):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid product status",
			})

		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})

		case errors.Is(err, repository.ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func (h *ProductHandler) UpdateProductStatus(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req request.UpdateProductStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	status := model.ProductStatus(req.Status)

	updatedProduct, err := h.service.UpdateProductStatus(ctx, id, status)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProductStatus):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, updatedProduct)

}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.service.DeleteProduct(ctx, id)

	if errors.Is(err, repository.ErrProductNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) StockIn(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req request.StockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	request := model.StockRequest{
		Quantity: req.Quantity,
	}

	updatedProduct, err := h.service.StockIn(ctx, id, request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProductQuantity):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid product quantity",
			})

		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, response.NewStockResponse(updatedProduct))
}

func (h *ProductHandler) StockOut(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req request.StockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	request := model.StockRequest{
		Quantity: req.Quantity,
	}

	updatedProduct, err := h.service.StockOut(ctx, id, request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProductQuantity):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid product quantity",
			})

		case errors.Is(err, repository.ErrInsufficientStock):
			c.JSON(http.StatusConflict, gin.H{
				"error": "Insufficient stock",
			})

		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, response.NewStockResponse(updatedProduct))
}
