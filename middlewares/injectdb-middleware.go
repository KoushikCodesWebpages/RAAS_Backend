package middleware

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InjectDB(client *mongo.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Inject the MongoDB client into the context
        c.Set("db", client)
        c.Next()
    }
}
