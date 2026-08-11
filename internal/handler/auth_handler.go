package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olivercruznaguit/inventory-system/internal/handler/request"
	"github.com/olivercruznaguit/inventory-system/internal/handler/response"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
	"github.com/olivercruznaguit/inventory-system/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// godoc Login
// @Summary      authentication for inventory system
// @Description  create a jwt token for session
// @tags  	     auth
// @Accept       json
// @Produce      json
// @Param        auth	  body      request.AuthRequest  true  "Login Request"
// @Success      200  {object}  response.AuthResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /auth/login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body"))
		return
	}

	token, err := ah.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmailAddress):
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid email address"))
		case errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, service.ErrIncorrectPassword):
			c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Incorrect password or email address"))
		default:
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{Token: token})
}
