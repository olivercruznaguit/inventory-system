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

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new category with the provided name
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category body request.CreateCategoryRequest true "Category details"
// @Success      201  {object}  response.CategoryResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      409  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /categories [post]
func (ch *CategoryHandler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body"))
		return
	}

	category := model.Category{
		Name: req.Name,
	}

	createdCategory, err := ch.service.CreateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryAlreadyExists) {
			c.JSON(http.StatusConflict, response.NewErrorResponse("Category already exists"))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusCreated, response.NewCategoryResponse(createdCategory))
}

// GetCategories godoc
// @Summary      Retrieve all categories
// @Description  Get a list of all categories
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  []response.CategoryResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /categories [get]
func (ch *CategoryHandler) GetCategories(c *gin.Context) {
	ctx := c.Request.Context()

	categories, err := ch.service.GetCategories(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		return
	}

	categoryResponses := make([]response.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		categoryResponses = append(categoryResponses, response.NewCategoryResponse(category))
	}

	c.JSON(http.StatusOK, categoryResponses)
}

// GetCategoryByID godoc
// @Summary 	   Retrieve a category by ID
// @Description    Get a category by its unique ID
// @Tags 		   Categories
// @Accept         json
// @Produce        json
// @param          id path int true "Category ID"
// @Success        200 {object} response.CategoryResponse
// @Failure        400 {object} response.ErrorResponse
// @Failure        404 {object} response.ErrorResponse
// @Failure        500 {object} response.ErrorResponse
// @Router         /categories/{id} [get]
func (ch *CategoryHandler) GetCategoryByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid category ID"))
		return
	}

	category, err := ch.service.GetCategoryByID(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Category not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		return
	}

	c.JSON(http.StatusOK, response.NewCategoryResponse(category))
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Update an existing category by its unique ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Category ID"
// @Param        category body      request.UpdateCategoryRequest  true  "Updated category details"
// @Success      200  {object}  response.CategoryResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      409  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /categories/{id} [put]
func (ch *CategoryHandler) UpdateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid category ID"))
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body"))
		return
	}

	category := model.Category{
		ID:   uint(id),
		Name: req.Name,
	}

	updatedCategory, err := ch.service.UpdateCategory(ctx, category)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrCategoryAlreadyExists):
			c.JSON(http.StatusConflict, response.NewErrorResponse("Category already exists"))

		case errors.Is(err, repository.ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Category not found"))

		default:
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		}

		return
	}

	c.JSON(http.StatusOK, response.NewCategoryResponse(updatedCategory))
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Delete an existing category by its unique ID
// @Tags         Categories
// @Param        id path int true "Category ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /categories/{id} [delete]
func (ch *CategoryHandler) DeleteCategory(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid category ID"))
		return
	}

	err = ch.service.DeleteCategory(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Category not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		return
	}

	c.Status(http.StatusNoContent)
}
