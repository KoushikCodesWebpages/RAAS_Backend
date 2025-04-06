package routes

import (
	"RAAS/config"
	"RAAS/handlers"
	"RAAS/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Main entry point to register routes
func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	r.Use(middleware.InjectDB(db))

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	// Auth routes
	r.POST("/signup", func(c *gin.Context) {
		handlers.SeekerSignUp(c)
	})
	r.POST("/login", func(c *gin.Context) {
		handlers.Login(c)
	})

	// Profile routes (newly added)
	profileHandler := handlers.NewProfileHandler(db)

	// Use the AuthMiddleware to protect these routes
	r.Use(middleware.AuthMiddleware(cfg))

	// Profile endpoints
	r.GET("/profile", profileHandler.RetrieveProfile)   // Retrieve Profile
	r.PUT("/profile", profileHandler.UpdateProfile)     // Update Profile
	r.PATCH("/profile", profileHandler.PatchProfile)    // Partial Update Profile
	r.DELETE("/profile", profileHandler.DeleteProfile)  // Delete Profile
}
