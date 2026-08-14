package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivercruznaguit/inventory-system/internal/auth"
	"github.com/olivercruznaguit/inventory-system/internal/handler/response"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type AuthMiddleware struct {
	parser TokenParser
}

type TokenParser interface {
	ParseToken(token string) (auth.CustomClaims, error)
}

func NewAuthMiddleware(parser TokenParser) *AuthMiddleware {
	return &AuthMiddleware{
		parser: parser,
	}
}

func (m *AuthMiddleware) Authenticate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid authorization header"))
		c.Abort()
		return
	}

	claims, err := m.parser.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid token"))
		c.Abort()
		return
	}

	c.Set("claims", claims)
	c.Next()
}

func (m *AuthMiddleware) RequireAdmin(c *gin.Context) {
	value, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Claims do not exist"))
		c.Abort()
		return
	}

	claims, ok := value.(auth.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid claims"))
		c.Abort()
		return
	}

	if claims.Role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, response.NewErrorResponse("Insufficient permissions"))
		c.Abort()
		return
	}

	c.Next()
}
