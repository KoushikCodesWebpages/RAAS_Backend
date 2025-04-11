package routes

import (
	"RAAS/config"
	"RAAS/handlers"
	"RAAS/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	//"github.com/gin-contrib/cors"
	"time"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {

	// Middleware: IP Whitelisting
	// r.Use(middleware.IPWhitelistMiddleware([]string{
	// 	"136.232.10.146", // Koushik IP
	// 	"5.6.7.8",        // Frontend dev's IP
	// }))

	// // Middleware: CORS Configuration
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"}, // Update as needed
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// Middleware: Inject DB into context
	r.Use(middleware.InjectDB(db))

	// Health check route (now serving an HTML page)
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>RAAS Backend</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						text-align: center;
						margin-top: 50px;
						background-color: #f0f0f0;
					}
					h1 {
						color: #333;
					}
					p {
						color: #666;
					}
					.container {
						max-width: 600px;
						margin: 0 auto;
						padding: 20px;
						background-color: white;
						border-radius: 8px;
						box-shadow: 0 2px 5px rgba(0,0,0,0.1);
					}
				</style>
			</head>
			<body>
				<div class="container">
					<h1>Welcome to RAAS Backend</h1>
					<p>This is the backend server for the RAAS application, built with Gin.</p>
					<p><strong>Status:</strong> Up and Running</p>
					<p><strong>Current Time:</strong> %s</p>
					<p>Explore our API endpoints for authentication, profiles, and job data!</p>
				</div>
			</body>
			</html>
		`, time.Now().Format("2006-01-02 15:04:05"))
	})

	r.POST("/api/reset-db", handlers.ResetDBHandler)
	// AUTH ROUTES
	r.POST("/signup", handlers.SeekerSignUp)
	r.GET("/verify-email", handlers.VerifyEmail)
	r.POST("/login", handlers.Login)

	// PROFILE routes
	profileHandler := handlers.NewProfileHandler(db)
	profileRoutes := r.Group("/profile")
	profileRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		profileRoutes.GET("", profileHandler.RetrieveProfile)   // Retrieve Profile
		profileRoutes.PUT("", profileHandler.UpdateProfile)     // Update Profile
		profileRoutes.PATCH("", profileHandler.PatchProfile)    // Partial Update Profile
		profileRoutes.DELETE("", profileHandler.DeleteProfile)  // Delete Profile
	}

	//PERSONAL INFO routes

	personalInfoRoutes:=r.Group("/personal-info")
	personalInfoRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		personalInfoRoutes.POST("", handlers.CreatePersonalInfo)
		personalInfoRoutes.GET("", handlers.GetPersonalInfo)
		personalInfoRoutes.PUT("", handlers.UpdatePersonalInfo)
		personalInfoRoutes.PATCH("", handlers.PatchPersonalInfo)
	}

	//PROFESSIONAL SUMMARY routes

	professionalSummaryRoutes := r.Group("/professional-summary")
	professionalSummaryRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		professionalSummaryRoutes.POST("", handlers.CreateProfessionalSummary)
		professionalSummaryRoutes.GET("", handlers.GetProfessionalSummary)
		professionalSummaryRoutes.PUT("", handlers.UpdateProfessionalSummary)
	}

	//WORKEXPERIENCE routes

	workExpRoutes := r.Group("/work-experience")
	workExpRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		workExpRoutes.POST("", handlers.CreateWorkExperience)
		workExpRoutes.GET("", handlers.GetWorkExperience)
		workExpRoutes.PATCH("/:id", handlers.PatchWorkExperience)
		workExpRoutes.DELETE("/:id", handlers.DeleteWorkExperience)
	}

	//EDUCATION
	educationRoutes := r.Group("/education")
	educationRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		educationRoutes.POST("", handlers.CreateEducation)
		educationRoutes.GET("", handlers.GetEducation)
		educationRoutes.PUT("/:id", handlers.PutEducation)
		educationRoutes.DELETE("/:id", handlers.DeleteEducation)
	}

	//CERTIFICATES
	certificateRoutes := r.Group("/certificates")
	certificateRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		certificateRoutes.POST("", handlers.CreateCertificate)
		certificateRoutes.GET("", handlers.GetCertificates)
		certificateRoutes.PUT("/:id", handlers.PutCertificate)
		certificateRoutes.DELETE("/:id", handlers.DeleteCertificate)
	}

	//LANGUAGE
	languageRoutes := r.Group("/languages")
	languageRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		languageRoutes.POST("", handlers.CreateLanguage)
		languageRoutes.GET("", handlers.GetLanguages)
		languageRoutes.PUT("/:id", handlers.PutLanguage)
		languageRoutes.DELETE("/:id", handlers.DeleteLanguage)
	}

	// JOB TITLES routes
	jobTitlesRoutes := r.Group("/jobtitles")
	jobTitlesRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		jobTitlesRoutes.POST("", handlers.CreateJobTitle)
		jobTitlesRoutes.GET("", handlers.GetJobTitle)
		jobTitlesRoutes.PUT("", handlers.UpdateJobTitle)
		jobTitlesRoutes.PATCH("", handlers.PatchJobTitle)
	}

	// JOB DATA ROUTES
	jobRetrievalRoutes := r.Group("/api/jobs")
	jobRetrievalRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		jobRetrievalRoutes.GET("", handlers.JobRetrievalHandler)
	}

	// MATCH SCORE ROUTES
	matchScoreRoutes := r.Group("/matchscore")
	matchScoreRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		matchScoreRoutes.POST("", handlers.MatchScorePOST)
		matchScoreRoutes.GET("", handlers.MatchScoreGET)
	}
	




	

	

}