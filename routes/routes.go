package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"RAAS/handlers"
	"RAAS/models"
)

// SetupRoutes - Registers all routes
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Health check route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running"})
	})

	// Auth routes
	r.POST("/signup", func(c *gin.Context) {
		handlers.SeekerSignUp(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		handlers.Login(c)
	})
	
	r.GET("/verify-email", func(c *gin.Context) {
		c.Set("db", db)
		handlers.VerifyEmailHandler(c)
	})
	
	// CRUD for LinkedinJobMetadata
	SetupGenericRoutes[models.LinkedinJobMetadata](r, db, "/linkedin-job-metadata")
}

func SetupGenericRoutes[T any](r *gin.Engine, db *gorm.DB, baseRoute string) {
	handler := handlers.NewGenericHandler[T](db)

	group := r.Group(baseRoute)
	{
		group.POST("/", handler.Create)
		group.POST("/bulk", handler.BulkCreate)
		group.POST("/upload-csv", handler.UploadCSV)
		group.GET("/:id", handler.GetByID)
		group.GET("/", handler.GetAll)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	}
}
