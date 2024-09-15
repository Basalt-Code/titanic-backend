package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	CustomClaims
}

type CustomClaims struct {
	UserID    string
	SessionID string
	Role      string
}

type RefreshToken struct {
	Token     string
	ExpiredAt time.Time
}
