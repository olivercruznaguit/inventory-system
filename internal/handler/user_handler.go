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

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// godoc Register
// @Summary      Create user for inventory system
// @Description  Register a new user with email and password
// @tags  	     auth
// @Accept       json
// @Produce      json
// @Param        auth 	body      request.AuthRequest  true  "Register Request"
// @Success      201  	{object}  response.UserResponse
// @Failure      400  	{object}  response.ErrorResponse
// @Failure      409  	{object}  response.ErrorResponse
// @Failure      500  	{object}  response.ErrorResponse
// @Router       /auth/register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body"))
		return
	}

	user, err := uh.service.Register(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmailAddress):
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid email address"))
		case errors.Is(err, repository.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, response.NewErrorResponse("Email already exist"))
		default:
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Internal server error"))
		}
		return
	}

	c.JSON(http.StatusCreated, response.NewUserResponse(user))
}
