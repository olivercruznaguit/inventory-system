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

	if filterStatusStr != string(model.ProductStatusInactive) && filterStatusStr != string(model.ProductStatusActive) && filterStatusStr != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

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

	pagination := model.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	var filter = model.ProductFilter{}
	filter.Status = model.ProductStatus(filterStatusStr)
	filter.Pagination = pagination

	products, err := h.service.GetProducts(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productResponses []response.ProductResponse
	for _, product := range products.Products {
		productResponses = append(productResponses, response.ProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			Price:  product.Price,
			Status: string(product.Status),
		})
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
	if err != nil {
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

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
		return
	}

	product := model.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	createdProduct, err := h.service.CreateProduct(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product := model.Product{
		ID:    uint(id),
		Name:  req.Name,
		Price: req.Price,
	}

	updatedProduct, err := h.service.UpdateProduct(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
