package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func InjectAuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Just for testing: extract from header or hardcode
		// Replace with actual auth logic
		authHeader := c.GetHeader("X-Auth-User-ID")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing auth user ID"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(authHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auth user ID"})
			c.Abort()
			return
		}

		c.Set("authUserID", userID)
		c.Next()
	}
}
