package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivercruznaguit/inventory-system/internal/auth"
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization header",
		})
		c.Abort()
		return
	}

	claims, err := m.parser.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("claims", claims)
	c.Next()
}
