package routes

import (
	"RAAS/handlers/features"
	"RAAS/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"RAAS/config"
)

func SetupFeatureRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// PROFILE routes
	profileHandler := features.NewProfileHandler(db)
	profileRoutes := r.Group("/profile")
	profileRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		profileRoutes.GET("", profileHandler.RetrieveProfile)   // Retrieve Profile
		profileRoutes.PUT("", profileHandler.UpdateProfile)     // Update Profile
		profileRoutes.PATCH("", profileHandler.PatchProfile)    // Partial Update Profile
		profileRoutes.DELETE("", profileHandler.DeleteProfile)  // Delete Profile
	}

	// JOB DATA routes
	jobRetrievalRoutes := r.Group("/api/jobs")
	jobRetrievalRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		jobRetrievalRoutes.GET("", features.JobRetrievalHandler)
	}

	// JOB METADATA routes
	jobDataHandler := features.NewJobDataHandler(db)
	jobMetaRoutes := r.Group("/api/job-data")
	jobMetaRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		jobMetaRoutes.GET("/linkedin", jobDataHandler.GetLinkedInJobs)
		jobMetaRoutes.GET("/xing", jobDataHandler.GetXingJobs)
		jobMetaRoutes.GET("/linkedin/links", jobDataHandler.GetLinkedInLinks)
		jobMetaRoutes.GET("/xing/links", jobDataHandler.GetXingLinks)
		jobMetaRoutes.GET("/linkedin/descriptions", jobDataHandler.GetLinkedInDescriptions)
		jobMetaRoutes.GET("/xing/descriptions", jobDataHandler.GetXingDescriptions)
	}

}

