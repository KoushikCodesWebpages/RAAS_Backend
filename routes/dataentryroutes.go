package routes

// import (
// 	"RAAS/config"
// 	"RAAS/handlers/dataentry"
// 	"RAAS/handlers/features"
// 	"RAAS/middlewares"
// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func SetupDataEntryRoutes(r *gin.Engine, client *mongo.Client, cfg *config.Config) {

// 	// TIMELINE
// 	timeline := r.Group("/user/entry-progress")
// 	timeline.Use(middleware.AuthMiddleware()) // Middleware to authenticate JWT

// 	// Define the route for getting the next entry step
// 	timeline.GET("", features.GetNextEntryStep(client))

// 	// PERSONAL INFO routes
// 	personalInfoHandler := dataentry.NewPersonalInfoHandler(client)
// 	personalInfoRoutes := r.Group("/personal-info")
// 	personalInfoRoutes.Use(middleware.AuthMiddleware())
// 	{
// 		personalInfoRoutes.POST("", personalInfoHandler.CreatePersonalInfo)  // Create Personal Info
// 		personalInfoRoutes.GET("", personalInfoHandler.GetPersonalInfo)      // Get Personal Info
// 		personalInfoRoutes.PUT("", personalInfoHandler.UpdatePersonalInfo)   // Update Personal Info
// 		personalInfoRoutes.PATCH("", personalInfoHandler.PatchPersonalInfo)  // Partial Update Personal Info
// 	}

// 	// PROFESSIONAL SUMMARY routes
// 	professionalSummaryHandler := dataentry.NewProfessionalSummaryHandler(client)
// 	professionalSummaryRoutes := r.Group("/professional-summary")
// 	professionalSummaryRoutes.Use(middleware.AuthMiddleware())
// 	{
// 		professionalSummaryRoutes.POST("", professionalSummaryHandler.CreateProfessionalSummary)
// 		professionalSummaryRoutes.GET("", professionalSummaryHandler.GetProfessionalSummary)
// 		professionalSummaryRoutes.PUT("", professionalSummaryHandler.UpdateProfessionalSummary)
// 	}

// 	// CERTIFICATES routes
// 	certificateHandler := dataentry.NewCertificateHandler(client)
// 	certificateRoutes := r.Group("/certificates")
// 	certificateRoutes.Use(middleware.AuthMiddleware())
// 	{
// 		certificateRoutes.POST("", certificateHandler.CreateCertificate)
// 		certificateRoutes.GET("", certificateHandler.GetCertificates)
// 		certificateRoutes.PUT(":id", certificateHandler.PatchCertificate)
// 		certificateRoutes.DELETE(":id", certificateHandler.DeleteCertificate)
// 	}

// 	// LANGUAGES routes	
// 	languageHandler := dataentry.NewLanguageHandler(client)
// 	languageRoutes := r.Group("/languages")
// 	languageRoutes.Use(middleware.AuthMiddleware())
// 	{
// 		languageRoutes.POST("", languageHandler.CreateLanguage)
// 		languageRoutes.GET("", languageHandler.GetLanguages)
// 		languageRoutes.PUT(":id", languageHandler.PatchLanguage)
// 		languageRoutes.DELETE(":id", languageHandler.DeleteLanguage)
// 	}

// 	// JOB TITLES routes
// 	jobTitleHandler := dataentry.NewJobTitleHandler(client)
// 	jobTitleRoutes := r.Group("/jobtitles")
// 	jobTitleRoutes.Use(middleware.AuthMiddleware())
// 	{
// 		jobTitleRoutes.POST("", jobTitleHandler.CreateJobTitle)
// 		jobTitleRoutes.GET("", jobTitleHandler.GetJobTitle)
// 		jobTitleRoutes.PUT("", jobTitleHandler.UpdateJobTitle)
// 		jobTitleRoutes.PATCH("", jobTitleHandler.PatchJobTitle)
// 	}
// }
