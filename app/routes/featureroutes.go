package routes

import (
	"RAAS/core/config"
	"RAAS/core/middlewares"
	"RAAS/internal/handlers/features/generation"
	"RAAS/internal/handlers/features/jobs"
	"RAAS/internal/handlers/features/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupFeatureRoutes(r *gin.Engine, client *mongo.Client, cfg *config.Config) {
	// Inject MongoDB into context
	r.Use(middleware.InjectDB(client))

	// Auth Middleware + Pagination helpers
	auth := middleware.AuthMiddleware()
	paginate := middleware.PaginationMiddleware

	// === USER ===

	seekerProfileHandler := user.NewSeekerProfileHandler()
	r.Group("/profile", auth).
		GET("", seekerProfileHandler.GetSeekerProfile)

	savedJobsHandler := user.NewSavedJobsHandler()
	r.Group("/saved-jobs", auth, paginate).
		POST("", savedJobsHandler.SaveJob).
		GET("", savedJobsHandler.GetSavedJobs)

	selectedJobsHandler := user.NewSelectedJobsHandler()
	r.Group("/api/selected-jobs", auth, paginate).
		POST("", selectedJobsHandler.PostSelectedJob).
		GET("", selectedJobsHandler.GetSelectedJobs)

	// === JOBS ===

	r.Group("/api/jobs", auth, paginate).
		GET("", jobs.JobRetrievalHandler)

	linkProviderHandler := jobs.NewLinkProviderHandler()
	r.Group("/provide-link", auth).
		POST("", linkProviderHandler.PostAndGetLink)

	// === GENERATION ===

	coverLetterHandler := generation.NewCoverLetterHandler()
	r.Group("/generate-cover-letter", auth).
		POST("", coverLetterHandler.PostCoverLetter)

	resumeHandler := generation.NewResumeHandler()
	r.Group("/generate-resume", auth).
		POST("", resumeHandler.PostResume)

		

	// // JOB METADATA routes
	// jobDataHandler := features.NewJobDataHandler()
	// jobMetaRoutes := r.Group("/api/job-data")
	// jobMetaRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	jobMetaRoutes.GET("", jobDataHandler.GetAllJobs)
	// }


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


}
