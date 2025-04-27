package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// PaginationMiddleware handles pagination logic
func PaginationMiddleware(c *gin.Context) {
	// Get the "page" and "limit" query parameters
	page := c.DefaultQuery("page", "1") // Default to 1 if not provided
	limit := c.DefaultQuery("limit", "10") // Default to 10 if not provided

	// Convert page and limit to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1 // Default to page 1 if invalid
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 10 // Default to limit of 10 if invalid
	}

	// Store pagination information in context
	c.Set("pagination", gin.H{
		"page":  pageInt,
		"limit": limitInt,
	})

	// Proceed to the next handler
	c.Next()
}
