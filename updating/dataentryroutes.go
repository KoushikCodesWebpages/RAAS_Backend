package routes

import (
	"RAAS/core/config"
	"RAAS/core/middlewares"

	"RAAS/internal/handlers/features/user"
	"RAAS/internal/handlers/preference"


	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupDataEntryRoutes(r *gin.Engine, client *mongo.Client, cfg *config.Config) {
	r.Use(middleware.InjectDB(client))
	// TIMELINE
	timeline := r.Group("/user/entry-progress/check")
	timeline.Use(middleware.AuthMiddleware()) // Middleware to authenticate JWT

	// Define the route for getting the next entry step
	timeline.GET("", user.GetNextEntryStep())

	// PERSONAL INFO routes
	personalInfoHandler := preference.NewPersonalInfoHandler()
	personalInfoRoutes := r.Group("/personal-info")
	personalInfoRoutes.Use(middleware.AuthMiddleware())
	{
		personalInfoRoutes.POST("", personalInfoHandler.CreatePersonalInfo)
		personalInfoRoutes.GET("", personalInfoHandler.GetPersonalInfo)    
		personalInfoRoutes.PUT("", personalInfoHandler.UpdatePersonalInfo)   
		personalInfoRoutes.PATCH("", personalInfoHandler.PatchPersonalInfo)  
	}

	// PROFESSIONAL SUMMARY routes
	professionalSummaryHandler := preference.NewProfessionalSummaryHandler()
	professionalSummaryRoutes := r.Group("/professional-summary")
	professionalSummaryRoutes.Use(middleware.AuthMiddleware())
	{
		professionalSummaryRoutes.POST("", professionalSummaryHandler.CreateProfessionalSummary)
		professionalSummaryRoutes.GET("", professionalSummaryHandler.GetProfessionalSummary)
		professionalSummaryRoutes.PUT("", professionalSummaryHandler.UpdateProfessionalSummary)
	}


	workExperienceHandler := preference.NewWorkExperienceHandler()
	workExperienceRoutes := r.Group("/work-experience")
	workExperienceRoutes.Use(middleware.AuthMiddleware())
	{
		workExperienceRoutes.POST("", workExperienceHandler.CreateWorkExperience)
		workExperienceRoutes.GET("", workExperienceHandler.GetWorkExperience)
		// workExperienceRoutes.PUT("", workExperienceHandler.UpdateWorkExperience)
	}


	educationHandler := preference.NewEducationHandler()
	educationRoutes := r.Group("/education")
	educationRoutes.Use(middleware.AuthMiddleware())
	{
		educationRoutes.POST("", educationHandler.CreateEducation)
		educationRoutes.GET("", educationHandler.GetEducation)
		// educationRoutes.PUT("", educationHandler.UpdateEducation)
	}


	// CERTIFICATES routes
	certificateHandler := preference.NewCertificateHandler()
	certificateRoutes := r.Group("/certificates")
	certificateRoutes.Use(middleware.AuthMiddleware())
	{
		certificateRoutes.POST("", certificateHandler.CreateCertificate)
		certificateRoutes.GET("", certificateHandler.GetCertificates)
		// certificateRoutes.PUT(":id", certificateHandler.PatchCertificate)
		// certificateRoutes.DELETE(":id", certificateHandler.DeleteCertificate)
	}

	// LANGUAGES routes	
	languageHandler := preference.NewLanguageHandler()
	languageRoutes := r.Group("/languages")
	languageRoutes.Use(middleware.AuthMiddleware())
	{
		languageRoutes.POST("", languageHandler.CreateLanguage)
		languageRoutes.GET("", languageHandler.GetLanguages)
		// languageRoutes.PUT(":id", languageHandler.PatchLanguage)
		// languageRoutes.DELETE(":id", languageHandler.DeleteLanguage)
	}

	// JOB TITLES routes
	jobTitleHandler := preference.NewJobTitleHandler()
	jobTitleRoutes := r.Group("/jobtitles")
	jobTitleRoutes.Use(middleware.AuthMiddleware())
	{
		jobTitleRoutes.POST("", jobTitleHandler.CreateJobTitleOnce)
		jobTitleRoutes.GET("", jobTitleHandler.GetJobTitle)
		// jobTitleRoutes.PATCH("", jobTitleHandler.PatchJobTitle)
	}
}
