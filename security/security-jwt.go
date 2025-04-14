package security

import (
	"RAAS/config"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strings"
	"time"
)

var jwtSecret = []byte(config.Cfg.JWTSecretKey)

// CustomClaims struct for JWT claims with uuid.UUID for UserID
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates an access token for the user
func GenerateJWT(userID uuid.UUID, email, role string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.Cfg.AccessTokenLifetime))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT checks the token string and returns custom claims
func ValidateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("could not parse JWT claims")
	}

	return claims, nil
}

// ParseJWTFromHeader extracts the token from Authorization header and validates it
func ParseJWTFromHeader(authHeader string) (*CustomClaims, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("missing Bearer token")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	return ValidateJWT(tokenStr)
}

// Helper functions for common checks
func IsRole(claims *CustomClaims, role string) bool {
	return claims.Role == role
}

func GetUserID(claims *CustomClaims) uuid.UUID {
	return claims.UserID
}

func GetEmail(claims *CustomClaims) string {
	return claims.Email
}