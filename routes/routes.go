package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"RAAS/handlers"
	//"RAAS/models"
	"RAAS/config"
	"RAAS/middlewares"
)

// SetupRoutes - Registers all routes
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.Use(middleware.InjectDB(db))
	// Health check route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running"})
	})

	// Auth routes
	r.POST("/signup", func(c *gin.Context) {
		handlers.SeekerSignUp(c)
	})

	// r.GET("/auth/google", handlers.GoogleLoginHandler)
	// r.GET("/auth/google/callback", handlers.GoogleCallbackHandler)
	r.GET("/verify-email", func(c *gin.Context) {
		c.Set("db", db)
		handlers.VerifyEmailHandler(c)
	})

	r.POST("/login", func(c *gin.Context) {
		handlers.Login(c)
	})
	
	//SetupProtectedRoutes[models.Product](router, db, "/products")

	
	// CRUD for LinkedinJobMetadata
}

// routes/protected_routes.go


func SetupProtectedRoutes[T any](r *gin.Engine, db *gorm.DB, baseRoute string, cfg *config.Config) {
	handler := handlers.NewProtectedHandler[T](db)

	group := r.Group(baseRoute, middleware.AuthMiddleware(cfg)) // middleware uses cfg now
	{
		group.POST("/", handler.CreateProtected)
		group.POST("/bulk", handler.BulkCreateProtected)
		group.GET("/:id", handler.GetByIDProtected)
		group.GET("/", handler.GetAllProtected)
		group.PUT("/:id", handler.UpdateProtected)
		group.DELETE("/:id", handler.DeleteProtected)
	}
}

