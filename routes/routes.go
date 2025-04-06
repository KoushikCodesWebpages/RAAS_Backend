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
	// Use InjectDB middleware to make the database instance available in the context
	r.Use(middleware.InjectDB(db))

	// Health check route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	// Auth routes
	r.POST("/signup", handlers.SeekerSignUp)
	r.POST("/login", handlers.Login)

	// Profile routes
	profileHandler := handlers.NewProfileHandler(db)
	// Protect profile routes with AuthMiddleware
	profileRoutes := r.Group("/profile")
	profileRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		profileRoutes.GET("", profileHandler.RetrieveProfile)   // Retrieve Profile
		profileRoutes.PUT("", profileHandler.UpdateProfile)     // Update Profile
		profileRoutes.PATCH("", profileHandler.PatchProfile)    // Partial Update Profile
		profileRoutes.DELETE("", profileHandler.DeleteProfile)  // Delete Profile
	}

	// Job titles routes
	jobTitlesRoutes := r.Group("/jobtitles")
	jobTitlesRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		jobTitlesRoutes.POST("", handlers.CreateJobTitle)
		jobTitlesRoutes.GET("", handlers.GetJobTitles)
		jobTitlesRoutes.PUT("", handlers.UpdateJobTitle)
		jobTitlesRoutes.PATCH("", handlers.PatchJobTitle)
	}
}
