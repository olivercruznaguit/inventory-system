package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type TokenService struct {
	secret []byte
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{
		secret: []byte(secret),
	}
}

func (ts *TokenService) GenerateToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "inventory-system",
		},
	})

	signedToken, err := token.SignedString(ts.secret)

	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return signedToken, nil
}

func (ts *TokenService) ParseToken(tokenString string) (CustomClaims, error) {
	var claims CustomClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return ts.secret, nil
		},
	)

	if err != nil {
		return CustomClaims{}, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return CustomClaims{}, errors.New("invalid token")
	}

	issuer, err := claims.GetIssuer()
	if err != nil || issuer != "inventory-system" {
		return CustomClaims{}, errors.New("invalid token issuer")
	}

	return claims, nil
}
