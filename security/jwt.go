package security

import (
	"RAAS/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"errors"
	"strings"
)

var appConfig *config.Config
var jwtSecret []byte

func init() {
	var err error
	appConfig, err = config.InitConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
	jwtSecret = []byte(appConfig.JWTSecretKey)
}

type CustomClaims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}


// GenerateJWT creates an access token with configured expiry
func GenerateJWT(userID uint, email, role string, cfg *config.Config) (string, error) {
    claims := CustomClaims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(cfg.AccessTokenLifetime))),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.JWTSecretKey))
}


// ValidateJWT verifies the token and returns claims if valid
func ValidateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ParseJWTFromHeader(authHeader string) (*CustomClaims, error) {
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}
	return claims, nil
}
