package routes

import (
	"RAAS/config"
	"RAAS/handlers"
	"RAAS/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/gin-contrib/cors"
	"time"
)

// Main entry point to register routes
func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {

	r.Use(middleware.IPWhitelistMiddleware([]string{
		//run thus ( curl ifconfig.me )

		"136.232.10.146", // Koushik IP
		"5.6.7.8", // Frontend dev's IP
	}))
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // update as needed
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Use InjectDB middleware to make the database instance available in the context
	r.Use(middleware.InjectDB(db))

	// Health check route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	// AUTH ROUTES
	r.POST("/signup", handlers.SeekerSignUp)
	r.GET("/verify-email", handlers.VerifyEmail)
	r.POST("/login", handlers.Login)

	// PROFILE routes
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

	// JOB TITLES routes
	jobTitlesRoutes := r.Group("/jobtitles")
	jobTitlesRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		jobTitlesRoutes.POST("", handlers.CreateJobTitle)
		jobTitlesRoutes.GET("", handlers.GetJobTitle)
		jobTitlesRoutes.PUT("", handlers.UpdateJobTitle)
		jobTitlesRoutes.PATCH("", handlers.PatchJobTitle)
	}


	//JOB DATA ROUTES

	r.GET("/api/jobs", handlers.JobRetrievalHandler)
}
