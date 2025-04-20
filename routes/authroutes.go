package routes

import (
	"RAAS/handlers/auth"
	"RAAS/config"
	
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine, cfg *config.Config) {
	// AUTH ROUTES
	r.POST("/signup", auth.SeekerSignUp)
	r.GET("/verify-email", auth.VerifyEmail)
	r.POST("/login", auth.Login)
	r.POST("/auth/forgot-password", auth.ForgotPasswordHandler)
	r.POST("/auth/admin-reset-token", auth.SystemInitiatedResetTokenHandler) // no email
	
	// Define the route for the reset password page
	r.GET("/reset-password", auth.ResetPasswordPage) // This should match the route you're using for the password reset page

	// Define the route for the actual password reset submission
	r.POST("/reset-password", auth.ResetPasswordHandler)

		
}
