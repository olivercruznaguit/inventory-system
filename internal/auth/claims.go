package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type CustomClaims struct {
	UserID uint           `json:"user_id"`
	Email  string         `json:"email"`
	Role   model.UserRole `json:"role"`
	jwt.RegisteredClaims
}
