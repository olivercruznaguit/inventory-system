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

// GetProducts godoc
// @Summary      Get products
// @Description  Get a list of products with optional filtering and pagination
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        page     query     int    false  "Page number"
// @Param        pageSize query     int    false  "Number of items per page"
// @Param        status   query     string false  "Product status (ACTIVE or INACTIVE)"
// @Param        search   query     string false  "Search term"
// @Param        categoryId query   int    false  "Category ID"
// @Param        sortBy   query     string false  "Sort by field"
// @Param        sortOrder query    string false  "Sort order (DESC or ASC)"
// @Success      200  {object}  response.ProductListResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products [get]
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
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid page number"))
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid page size"))
		return
	}

	var categoryID *uint
	if categoryIdStr != "" {
		categoryIdParsed, err := strconv.Atoi(categoryIdStr)
		if err != nil || categoryIdParsed < 1 {
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid category ID"))
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
			c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		return
	}

	productResponses := make([]response.ProductResponse, 0, len(products.Products))
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

// GetProductByID godoc
// @Summary      Get product by ID
// @Description  Get detailed information of a product by its unique ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  response.ProductResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid product ID"))
		return
	}

	product, err := h.service.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Product not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.NewProductResponse(product))
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product with the provided details
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product body request.CreateProductRequest true "Product details"
// @Success      201  {object}  response.ProductResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body"))
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
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid product minimum stock"))
			return
		}

		if errors.Is(err, repository.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Category not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusCreated, response.NewProductResponse(createdProduct))
}

// UpdateProduct godoc
// @Summary      Update an existing product
// @Description  Update the details of an existing product by its unique ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Product ID"
// @Param        product body request.UpdateProductRequest true "Updated product details"
// @Success      200  {object}  response.ProductResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid product ID"))
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body"))
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
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid product minimum stock"))

		case errors.Is(err, service.ErrInvalidProductStatus):
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid product status"))

		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Product not found"))

		case errors.Is(err, repository.ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Category not found"))

		default:
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		}

		return
	}

	c.JSON(http.StatusOK, response.NewProductResponse(updatedProduct))
}

// UpdateProductStatus godoc
// @Summary      Update product status
// @Description  Update the status of an existing product by its unique ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Product ID"
// @Param        status  body      string true  "Updated product status"
// @Success      200  {object}  response.ProductResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products/{id}/status [patch]
func (h *ProductHandler) UpdateProductStatus(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid product ID"))
		return
	}

	var req request.UpdateProductStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body"))
		return
	}

	status := model.ProductStatus(req.Status)

	updatedProduct, err := h.service.UpdateProductStatus(ctx, id, status)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProductStatus):
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid product status"))

		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Product not found"))

		default:
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		}

		return
	}

	c.JSON(http.StatusOK, response.NewProductResponse(updatedProduct))

}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete an existing product by its unique ID
// @Tags         Products
// @Param        id path int true "Product ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid product ID"))
		return
	}

	err = h.service.DeleteProduct(ctx, id)

	if errors.Is(err, repository.ErrProductNotFound) {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Product not found"))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		return
	}

	c.Status(http.StatusNoContent)
}
