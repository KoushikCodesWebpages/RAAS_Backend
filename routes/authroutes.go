// routes/auth_routes.go
package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"RAAS/config"
	"RAAS/handlers/auth"
	"RAAS/middlewares"
)

func SetupAuthRoutes(r *gin.Engine, cfg *config.Config) {
	// Rate limiter configurations
	signupLimiter := middleware.RateLimiterMiddleware(5, time.Minute)
	loginLimiter := middleware.RateLimiterMiddleware(10, time.Minute)
	forgotPassLimiter := middleware.RateLimiterMiddleware(3, time.Minute)
	resetPassLimiter := middleware.RateLimiterMiddleware(3, time.Minute)
	verifyEmailLimiter := middleware.RateLimiterMiddleware(10, time.Minute)
	googleLoginLimiter := middleware.RateLimiterMiddleware(10, time.Minute)
	googleCallbackLimiter := middleware.RateLimiterMiddleware(20, time.Minute)

	authGroup := r.Group("/auth")
	{
		// Google OAuth (rate-limited)
		authGroup.GET("/google/login", googleLoginLimiter, auth.GoogleLoginHandler)
		authGroup.GET("/google/callback", googleCallbackLimiter, auth.GoogleCallbackHandler)

		// Standard auth routes (rate-limited where necessary)
		authGroup.POST("/signup", signupLimiter, auth.SeekerSignUp)
		authGroup.GET("/verify-email", verifyEmailLimiter, auth.VerifyEmail)
		authGroup.POST("/login", loginLimiter, auth.Login)
		authGroup.POST("/forgot-password", forgotPassLimiter, auth.ForgotPasswordHandler)
		authGroup.POST("/admin-reset-token", auth.SystemInitiatedResetTokenHandler) // No limiter
		authGroup.GET("/reset-password", auth.ResetPasswordPage)                     // Optional
		authGroup.POST("/reset-password", resetPassLimiter, auth.ResetPasswordHandler)
	}
}
