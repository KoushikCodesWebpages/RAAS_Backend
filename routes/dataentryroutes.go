package routes

import (
	"RAAS/handlers/dataentry"
	"RAAS/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"RAAS/config"
)

func SetupDataEntryRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// PERSONAL INFO routes
	personalInfoHandler := dataentry.NewPersonalInfoHandler(db)
	personalInfoRoutes := r.Group("/personal-info")
	personalInfoRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		personalInfoRoutes.POST("", personalInfoHandler.CreatePersonalInfo)  // Create Personal Info
		personalInfoRoutes.GET("", personalInfoHandler.GetPersonalInfo)      // Get Personal Info
		personalInfoRoutes.PUT("", personalInfoHandler.UpdatePersonalInfo)   // Update Personal Info
		personalInfoRoutes.PATCH("", personalInfoHandler.PatchPersonalInfo)  // Partial Update Personal Info
	}
	
	// PROFESSIONAL SUMMARY routes
	professionalSummaryHandler := dataentry.NewProfessionalSummaryHandler(db)
	professionalSummaryRoutes := r.Group("/professional-summary")
	professionalSummaryRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		professionalSummaryRoutes.POST("", professionalSummaryHandler.CreateProfessionalSummary)
		professionalSummaryRoutes.GET("", professionalSummaryHandler.GetProfessionalSummary)
		professionalSummaryRoutes.PUT("", professionalSummaryHandler.UpdateProfessionalSummary)
	}
	
	// WORK EXPERIENCE routes
	workExperienceHandler := dataentry.NewWorkExperienceHandler(db)
	workExperienceRoutes := r.Group("/work-experience")
	workExperienceRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		workExperienceRoutes.POST("", workExperienceHandler.CreateWorkExperience)
		workExperienceRoutes.GET("", workExperienceHandler.GetWorkExperience)
		workExperienceRoutes.PATCH(":id", workExperienceHandler.PatchWorkExperience)
		workExperienceRoutes.DELETE(":id", workExperienceHandler.DeleteWorkExperience)
	}

	// EDUCATION routes
	educationHandler := dataentry.NewEducationHandler(db)
	educationRoutes := r.Group("/education")
	educationRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		educationRoutes.POST("", educationHandler.CreateEducation)
		educationRoutes.GET("", educationHandler.GetEducation)
		educationRoutes.PUT(":id", educationHandler.PutEducation)
		educationRoutes.DELETE(":id", educationHandler.DeleteEducation)
	}

	// CERTIFICATES routes
	certificateHandler := dataentry.NewCertificateHandler(db)
	certificateRoutes := r.Group("/certificates")
	certificateRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		certificateRoutes.POST("", certificateHandler.CreateCertificate)
		certificateRoutes.GET("", certificateHandler.GetCertificates)
		certificateRoutes.PUT(":id", certificateHandler.PutCertificate)
		certificateRoutes.DELETE(":id", certificateHandler.DeleteCertificate)
	}

	// LANGUAGES routes
	languageHandler := dataentry.NewLanguageHandler(db)
	languageRoutes := r.Group("/languages")
	languageRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		languageRoutes.POST("", languageHandler.CreateLanguage)
		languageRoutes.GET("", languageHandler.GetLanguages)
		languageRoutes.PUT(":id", languageHandler.PutLanguage)
		languageRoutes.DELETE(":id", languageHandler.DeleteLanguage)
	}

	// JOB TITLES routes
	jobTitleHandler := dataentry.NewJobTitleHandler(db)
	jobTitleRoutes := r.Group("/jobtitles")
	jobTitleRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		jobTitleRoutes.POST("", jobTitleHandler.CreateJobTitle)
		jobTitleRoutes.GET("", jobTitleHandler.GetJobTitle)
		jobTitleRoutes.PUT("", jobTitleHandler.UpdateJobTitle)
		jobTitleRoutes.PATCH("", jobTitleHandler.PatchJobTitle)
	}
}
