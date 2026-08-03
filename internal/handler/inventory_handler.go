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

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service: service,
	}
}

func (h *InventoryHandler) StockIn(c *gin.Context) {
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

	stockReq := model.StockRequest{
		Quantity:  req.Quantity,
		Reason:    req.Reason,
		ProductID: id,
	}

	updatedProduct, err := h.service.StockIn(ctx, stockReq)
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

func (h *InventoryHandler) StockOut(c *gin.Context) {
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

	stockReq := model.StockRequest{
		Quantity:  req.Quantity,
		Reason:    req.Reason,
		ProductID: id,
	}

	updatedProduct, err := h.service.StockOut(ctx, stockReq)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInsufficientStock):
			c.JSON(http.StatusBadRequest, gin.H{
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

func (h *InventoryHandler) GetStockMovements(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	stockMovements, err := h.service.GetStockMovements(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stockResponses := make([]response.StockMovementResponse, 0, len(stockMovements))
	for _, stock := range stockMovements {
		stockResponses = append(stockResponses, response.NewStockMovementResponse(stock))
	}

	stockResponse := response.StockMovementListResponse{
		Data: stockResponses,
	}

	c.JSON(http.StatusOK, stockResponse)
}

func (h *InventoryHandler) GetInventoryDashboard(c *gin.Context) {
	ctx := c.Request.Context()

	inventoryDashboard, err := h.service.GetInventoryDashboard(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, response.NewInventoryDashboardResponse(inventoryDashboard))
}
