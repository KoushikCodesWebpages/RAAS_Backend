package routes

import (
	"RAAS/config"
	"RAAS/handlers/dataentry"
	"RAAS/handlers/features"
	"RAAS/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupDataEntryRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {


	//TIMELINE 

	timeline := r.Group("/user/entry-progress")
    timeline.Use(middleware.AuthMiddleware()) // Middleware to authenticate JWT

    // Define the route for getting the next entry step
    timeline.GET("", features.GetNextEntryStep(db))
	// PERSONAL INFO routes

	personalInfoHandler := dataentry.NewPersonalInfoHandler(db)
	personalInfoRoutes := r.Group("/personal-info")
	personalInfoRoutes.Use(middleware.AuthMiddleware())
	{
		personalInfoRoutes.POST("", personalInfoHandler.CreatePersonalInfo)  // Create Personal Info
		personalInfoRoutes.GET("", personalInfoHandler.GetPersonalInfo)      // Get Personal Info
		personalInfoRoutes.PUT("", personalInfoHandler.UpdatePersonalInfo)   // Update Personal Info
		personalInfoRoutes.PATCH("", personalInfoHandler.PatchPersonalInfo)  // Partial Update Personal Info
	}
	
	// PROFESSIONAL SUMMARY routes
	professionalSummaryHandler := dataentry.NewProfessionalSummaryHandler(db)
	professionalSummaryRoutes := r.Group("/professional-summary")
	professionalSummaryRoutes.Use(middleware.AuthMiddleware())
	{
		professionalSummaryRoutes.POST("", professionalSummaryHandler.CreateProfessionalSummary)
		professionalSummaryRoutes.GET("", professionalSummaryHandler.GetProfessionalSummary)
		professionalSummaryRoutes.PUT("", professionalSummaryHandler.UpdateProfessionalSummary)
	}
	
	// WORK EXPERIENCE routes
	workExperienceHandler := dataentry.NewWorkExperienceHandler(db)
	workExperienceRoutes := r.Group("/work-experience")
	workExperienceRoutes.Use(middleware.AuthMiddleware())
	{
		workExperienceRoutes.POST("", workExperienceHandler.CreateWorkExperience)
		workExperienceRoutes.GET("", workExperienceHandler.GetWorkExperiences)
		workExperienceRoutes.PATCH(":id", workExperienceHandler.PatchWorkExperience)
		workExperienceRoutes.DELETE(":id", workExperienceHandler.DeleteWorkExperience)
	}

	// EDUCATION routes
	educationHandler := dataentry.NewEducationHandler(db)
	educationRoutes := r.Group("/education")
	educationRoutes.Use(middleware.AuthMiddleware())
	{
		educationRoutes.POST("", educationHandler.CreateEducation)
		educationRoutes.GET("", educationHandler.GetEducations)
		educationRoutes.PUT(":id", educationHandler.PatchEducation)
		educationRoutes.DELETE(":id", educationHandler.DeleteEducation)
	}


	// CERTIFICATES routes
	certificateHandler := dataentry.NewCertificateHandler(db)
	certificateRoutes := r.Group("/certificates")
	certificateRoutes.Use(middleware.AuthMiddleware())
	{
		certificateRoutes.POST("", certificateHandler.CreateCertificate)
		certificateRoutes.GET("", certificateHandler.GetCertificates)
		certificateRoutes.PUT(":id", certificateHandler.PutCertificate)
		certificateRoutes.DELETE(":id", certificateHandler.DeleteCertificate)
	}

	// LANGUAGES routes
	languageHandler := dataentry.NewLanguageHandler(db)
	languageRoutes := r.Group("/languages")
	languageRoutes.Use(middleware.AuthMiddleware())
	{
		languageRoutes.POST("", languageHandler.CreateLanguage)
		languageRoutes.GET("", languageHandler.GetLanguages)
		languageRoutes.PUT(":id", languageHandler.UpdateLanguage)
		languageRoutes.DELETE(":id", languageHandler.DeleteLanguage)
	}

	// JOB TITLES routes
	jobTitleHandler := dataentry.NewJobTitleHandler(db)
	jobTitleRoutes := r.Group("/jobtitles")
	jobTitleRoutes.Use(middleware.AuthMiddleware())
	{
		jobTitleRoutes.POST("", jobTitleHandler.CreateJobTitle)
		jobTitleRoutes.GET("", jobTitleHandler.GetJobTitle)
		jobTitleRoutes.PUT("", jobTitleHandler.UpdateJobTitle)
		jobTitleRoutes.PATCH("", jobTitleHandler.PatchJobTitle)
	}
}
