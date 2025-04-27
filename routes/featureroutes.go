package routes

import (
	"RAAS/config"
	"RAAS/handlers/features"
	// "RAAS/handlers/repo"
	"RAAS/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupFeatureRoutes(r *gin.Engine, client *mongo.Client, cfg *config.Config) {

	r.Use(middleware.InjectDB(client))
	// PROFILE routes
	seekerProfileHandler := features.NewSeekerProfileHandler()
	seekerProfileRoutes := r.Group("/profile")
	seekerProfileRoutes.Use(middleware.AuthMiddleware())
	{
		seekerProfileRoutes.GET("", seekerProfileHandler.GetSeekerProfile)
	}
	jobRetrievalRoutes := r.Group("/api/jobs")
	jobRetrievalRoutes.Use(middleware.AuthMiddleware())      // Auth middleware for authentication
	jobRetrievalRoutes.Use(middleware.PaginationMiddleware) // Pagination middleware for pagination logic
	{
		jobRetrievalRoutes.GET("", features.JobRetrievalHandler)
	}

	// JOB METADATA routes
	jobDataHandler := features.NewJobDataHandler()
	jobMetaRoutes := r.Group("/api/job-data")
	jobMetaRoutes.Use(middleware.AuthMiddleware())
	{
		jobMetaRoutes.GET("", jobDataHandler.GetAllJobs)
	}

	// selectedJobsHandler := features.NewSelectedJobsHandler(client)
	// selectedJobsRoutes := r.Group("/selected-jobs")
	// selectedJobsRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	selectedJobsRoutes.POST("", selectedJobsHandler.PostSelectedJob)
	// 	selectedJobsRoutes.GET("", selectedJobsHandler.GetSelectedJobs)
	// 	selectedJobsRoutes.PUT(":id", selectedJobsHandler.UpdateSelectedJob)
	// 	selectedJobsRoutes.DELETE(":id", selectedJobsHandler.DeleteSelectedJob)
	// }

	// matchScoreHandler := features.MatchScoreHandler{Client: client}
	// // Define the route group for match scores
	// matchScoreRoutes := r.Group("/matchscores")
	// matchScoreRoutes.Use(middleware.AuthMiddleware()) // If you want to secure it with authentication
	// {
	// 	// Route to get all match scores
	// 	matchScoreRoutes.GET("", matchScoreHandler.GetAllMatchScores)
	// }

	// CoverLetterHandler := features.NewCoverLetterHandler(client, cfg)
	// // CoverLetter generation route (authenticated)
	// coverLetterRoutes := r.Group("/generate-cover-letter")
	// coverLetterRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	coverLetterRoutes.POST("", CoverLetterHandler.PostCoverLetter)
	// }

	// cvHandler := features.NewCVHandler(client, cfg)
	// // CV generation route (authenticated)
	// cvRoutes := r.Group("/generate-cv")
	// cvRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	cvRoutes.POST("", cvHandler.PostCV)
	// }

	// CVHandler := features.NewCVDownloadHandler(client) // assuming constructor exists like NewCVHandler(client)

	// downloadCVRoutes := r.Group("/download-cv")
	// downloadCVRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	downloadCVRoutes.POST("", CVHandler.DownloadCV)
	// }

	// cvMetaHandler := features.NewCVDownloadHandler(client)
	// cvMetaRoutes := r.Group("/get-cv")
	// cvMetaRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	cvRoutes.GET("", cvMetaHandler.GetCVMetadata)
	// }

	// LinkProviderHandler := features.NewLinkProviderHandler(client)

	// // Link provider route (authenticated)
	// linkProviderRoutes := r.Group("/provide-link")
	// linkProviderRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	linkProviderRoutes.POST("", LinkProviderHandler.PostAndGetLink)
	// }
}
