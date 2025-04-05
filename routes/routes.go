package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"RAAS/controllers"
	"RAAS/repositories"
	"RAAS/models"
)

// SetupRoutes - Registers all routes
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Home Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running"})
	})

	// Authentication Route
	r.POST("/signup", func(c *gin.Context) {
		controllers.SeekerSignUp(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
        c.Set("db", db) // Pass DB instance to the controller
        controllers.Login(c)
    })

	// Register CRUD routes for LinkedinJobMetadata
	SetupGenericRoutes[models.LinkedinJobMetadata](r, db, "/linkedin-job-metadata")
}

// SetupGenericRoutes - Registers generic CRUD routes for any model
func SetupGenericRoutes[T any](r *gin.Engine, db *gorm.DB, baseRoute string) {
	repo := repositories.NewGeneralRepository[T](db)
	controller := controllers.NewGeneralController[T](repo)

	group := r.Group(baseRoute)
	{
		group.POST("/", controller.Create)
		group.POST("/bulk", controller.BulkCreate)
		group.POST("/upload-csv", controller.UploadCSV)
		group.GET("/:id", controller.GetByID)
		group.GET("/", controller.GetAll)
		group.PUT("/:id", controller.Update)
		group.DELETE("/:id", controller.Delete)
	}
}
