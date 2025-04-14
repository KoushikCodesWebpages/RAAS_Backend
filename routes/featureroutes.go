package routes

import (
	"RAAS/handlers/features"
	"RAAS/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"RAAS/config"
	"os"
)

func SetupFeatureRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// PROFILE routes
	seekerProfileHandler := features.NewSeekerProfileHandler(db)
	seekerProfileRoutes := r.Group("/profile")
	seekerProfileRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		seekerProfileRoutes.GET("", seekerProfileHandler.GetSeekerProfile)
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

	selectedJobsHandler := features.NewSelectedJobsHandler(db)
	selectedJobsRoutes := r.Group("/selected-jobs")
	selectedJobsRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		selectedJobsRoutes.POST("", selectedJobsHandler.PostSelectedJob)
		selectedJobsRoutes.GET("", selectedJobsHandler.GetSelectedJobs)
		selectedJobsRoutes.PUT(":id", selectedJobsHandler.UpdateSelectedJob)
		selectedJobsRoutes.DELETE(":id", selectedJobsHandler.DeleteSelectedJob)
	}



	matchScoreHandler := features.MatchScoreHandler{DB: db}
	// Define the route group for match scores
	matchScoreRoutes := r.Group("/matchscores")
	matchScoreRoutes.Use(middleware.AuthMiddleware(cfg)) // If you want to secure it with authentication
	{
		// Route to get all match scores
		matchScoreRoutes.GET("", matchScoreHandler.GetAllMatchScores)
	}
	
	CLHandler := features.NewCoverLetterHandler(db,cfg)

	// Cover letter generation route (authenticated)
	coverLetterRoutes := r.Group("/generate-cover-letter")
	coverLetterRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		coverLetterRoutes.POST("", CLHandler.PostCoverLetter)
	}

	mediaUploadHandler := features.NewMediaUploadHandler(features.GetBlobServiceClient(), os.Getenv("AZURE_BLOB_CONTAINER"))
	mediaRoutes := r.Group("/media")
	mediaRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		mediaRoutes.POST("/upload", mediaUploadHandler.HandleUpload)
	}

}

