package security

import (
	"RAAS/config"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strings"
	"time"
	//"log"
)

var jwtSecret []byte

func getJWTSecret() []byte {
	if jwtSecret == nil {
		jwtSecret = []byte(config.Cfg.Project.JWTSecretKey)
	}
	return jwtSecret
}

// CustomClaims struct for JWT claims with uuid.UUID for UserID
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates an access token for the user
func GenerateJWT(userID uuid.UUID, email, role string) (string, error) {
    // log.Println("AccessTokenLifetime:", config.Cfg.AccessTokenLifetime)
    // log.Println("JWTSecretKey:", config.Cfg.JWTSecretKey)

    claims := CustomClaims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.Cfg.Project.AccessTokenLifetime))),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(getJWTSecret())
}

func ValidateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token has expired")
			}
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("token is malformed")
			}
		}
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
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

	tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
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