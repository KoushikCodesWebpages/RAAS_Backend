package routes

import (
	"RAAS/handlers/features"
	"RAAS/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"RAAS/config"
	//"os"
)

func SetupFeatureRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// PROFILE routes
	seekerProfileHandler := features.NewSeekerProfileHandler(db)
	seekerProfileRoutes := r.Group("/profile")
	seekerProfileRoutes.Use(middleware.AuthMiddleware())
	{
		seekerProfileRoutes.GET("", seekerProfileHandler.GetSeekerProfile)
	}
	
	

	// JOB DATA routes
	jobRetrievalRoutes := r.Group("/api/jobs")
	jobRetrievalRoutes.Use(middleware.AuthMiddleware())
	{
		jobRetrievalRoutes.GET("", features.JobRetrievalHandler)
	}

	// JOB METADATA routes
	jobDataHandler := features.NewJobDataHandler(db)
	jobMetaRoutes := r.Group("/api/job-data")
	jobMetaRoutes.Use(middleware.AuthMiddleware())
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
	selectedJobsRoutes.Use(middleware.AuthMiddleware())
	{
		selectedJobsRoutes.POST("", selectedJobsHandler.PostSelectedJob)
		selectedJobsRoutes.GET("", selectedJobsHandler.GetSelectedJobs)
		selectedJobsRoutes.PUT(":id", selectedJobsHandler.UpdateSelectedJob)
		selectedJobsRoutes.DELETE(":id", selectedJobsHandler.DeleteSelectedJob)
	}



	matchScoreHandler := features.MatchScoreHandler{DB: db}
	// Define the route group for match scores
	matchScoreRoutes := r.Group("/matchscores")
	matchScoreRoutes.Use(middleware.AuthMiddleware()) // If you want to secure it with authentication
	{
		// Route to get all match scores
		matchScoreRoutes.GET("", matchScoreHandler.GetAllMatchScores)
	}
	
	CoverLetterHandler := features.NewCoverLetterHandler(db, cfg)
	// CoverLetter  generation route (authenticated)
	coverLetterRoutes := r.Group("/generate-cover-letter")
	coverLetterRoutes.Use(middleware.AuthMiddleware())
	{
		coverLetterRoutes.POST("",CoverLetterHandler.PostCoverLetter)
	}

	cvHandler := features.NewCVHandler(db, cfg)
	// CV generation route (authenticated)
	cvRoutes := r.Group("/generate-cv")
	cvRoutes.Use(middleware.AuthMiddleware())
	{
		cvRoutes.POST("", cvHandler.PostCV)
	}

	LinkProviderHandler := features.NewLinkProviderHandler(db)

	// Link provider route (authenticated)
	linkProviderRoutes := r.Group("/provide-link")
	linkProviderRoutes.Use(middleware.AuthMiddleware())
	{
		linkProviderRoutes.POST("", LinkProviderHandler.PostAndGetLink)
	}


}

