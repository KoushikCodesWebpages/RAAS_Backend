package routes

import (
	"RAAS/config"
	"RAAS/handlers/dataentry"
	"RAAS/handlers/features"
	"RAAS/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupDataEntryRoutes(r *gin.Engine, client *mongo.Client, cfg *config.Config) {
	r.Use(middleware.InjectDB(client))
	// TIMELINE
	timeline := r.Group("/user/entry-progress/check")
	timeline.Use(middleware.AuthMiddleware()) // Middleware to authenticate JWT

	// Define the route for getting the next entry step
	timeline.GET("", features.GetNextEntryStep())

	// PERSONAL INFO routes
	personalInfoHandler := dataentry.NewPersonalInfoHandler()
	personalInfoRoutes := r.Group("/personal-info")
	personalInfoRoutes.Use(middleware.AuthMiddleware())
	{
		personalInfoRoutes.POST("", personalInfoHandler.CreatePersonalInfo)  // Create Personal Info
		personalInfoRoutes.GET("", personalInfoHandler.GetPersonalInfo)      // Get Personal Info
		personalInfoRoutes.PUT("", personalInfoHandler.UpdatePersonalInfo)   // Update Personal Info
		personalInfoRoutes.PATCH("", personalInfoHandler.PatchPersonalInfo)  // Partial Update Personal Info
	}

	// PROFESSIONAL SUMMARY routes
	professionalSummaryHandler := dataentry.NewProfessionalSummaryHandler()
	professionalSummaryRoutes := r.Group("/professional-summary")
	professionalSummaryRoutes.Use(middleware.AuthMiddleware())
	{
		professionalSummaryRoutes.POST("", professionalSummaryHandler.CreateProfessionalSummary)
		professionalSummaryRoutes.GET("", professionalSummaryHandler.GetProfessionalSummary)
		professionalSummaryRoutes.PUT("", professionalSummaryHandler.UpdateProfessionalSummary)
	}


	workExperienceHandler := dataentry.NewWorkExperienceHandler()
	workExperienceRoutes := r.Group("/work-experience")
	workExperienceRoutes.Use(middleware.AuthMiddleware())
	{
		workExperienceRoutes.POST("", workExperienceHandler.CreateWorkExperience)
		workExperienceRoutes.GET("", workExperienceHandler.GetWorkExperience)
		// workExperienceRoutes.PUT("", workExperienceHandler.UpdateWorkExperience)
	}


	educationHandler := dataentry.NewEducationHandler()
	educationRoutes := r.Group("/education")
	educationRoutes.Use(middleware.AuthMiddleware())
	{
		educationRoutes.POST("", educationHandler.CreateEducation)
		educationRoutes.GET("", educationHandler.GetEducation)
		// educationRoutes.PUT("", educationHandler.UpdateEducation)
	}


	// CERTIFICATES routes
	certificateHandler := dataentry.NewCertificateHandler()
	certificateRoutes := r.Group("/certificates")
	certificateRoutes.Use(middleware.AuthMiddleware())
	{
		certificateRoutes.POST("", certificateHandler.CreateCertificate)
		certificateRoutes.GET("", certificateHandler.GetCertificates)
		// certificateRoutes.PUT(":id", certificateHandler.PatchCertificate)
		// certificateRoutes.DELETE(":id", certificateHandler.DeleteCertificate)
	}

	// LANGUAGES routes	
	languageHandler := dataentry.NewLanguageHandler()
	languageRoutes := r.Group("/languages")
	languageRoutes.Use(middleware.AuthMiddleware())
	{
		languageRoutes.POST("", languageHandler.CreateLanguage)
		languageRoutes.GET("", languageHandler.GetLanguages)
		// languageRoutes.PUT(":id", languageHandler.PatchLanguage)
		// languageRoutes.DELETE(":id", languageHandler.DeleteLanguage)
	}

	// JOB TITLES routes
	jobTitleHandler := dataentry.NewJobTitleHandler()
	jobTitleRoutes := r.Group("/jobtitles")
	jobTitleRoutes.Use(middleware.AuthMiddleware())
	{
		jobTitleRoutes.POST("", jobTitleHandler.CreateJobTitleOnce)
		jobTitleRoutes.GET("", jobTitleHandler.GetJobTitle)
		// jobTitleRoutes.PATCH("", jobTitleHandler.PatchJobTitle)
	}
}
