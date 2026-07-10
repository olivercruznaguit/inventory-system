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

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (ch *CategoryHandler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	category := model.Category{
		Name: req.Name,
	}

	createdCategory, err := ch.service.CreateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Category already exist"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, response.NewCategoryResponse(createdCategory))
}

func (ch *CategoryHandler) GetCategories(c *gin.Context) {
	ctx := c.Request.Context()

	categories, err := ch.service.GetCategories(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	categoryResponses := make([]response.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		categoryResponses = append(categoryResponses, response.NewCategoryResponse(category))
	}

	c.JSON(http.StatusOK, categoryResponses)
}

func (ch *CategoryHandler) GetCategoryByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := ch.service.GetCategoryByID(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, response.NewCategoryResponse(category))
}
