package security

import (
    "github.com/golang-jwt/jwt/v4"
    "RAAS/config"
    "time"
	"log"
)

// Secret key (use env variable in production)
var jwtSecret = getJwtSecretKey()

// GenerateJWT creates a JWT token for the user
func GenerateJWT(email, role string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": email,
        "role":  role,
        "exp":   time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
    })

    return token.SignedString(jwtSecret)
}

// ValidateJWT checks the token validity
func ValidateJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
}


func getJwtSecretKey() string {
    config, err := config.InitConfig()
    if err != nil {
        log.Fatalf("Error initializing config: %v", err)
    }
    return config.JWTSecretKey 
}