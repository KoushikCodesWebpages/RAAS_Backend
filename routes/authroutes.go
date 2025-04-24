// routes/auth_routes.go
package routes

import (
	"RAAS/config"
	"RAAS/handlers/auth"
	"RAAS/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine, cfg *config.Config) {
	// Rate limit settings
	signupLimiter := middleware.RateLimiterMiddleware(5,time.Minute)
	loginLimiter := middleware.RateLimiterMiddleware(10,time.Minute)
	forgotPassLimiter := middleware.RateLimiterMiddleware(3,time.Minute)
	resetPassLimiter := middleware.RateLimiterMiddleware(3,time.Minute)
	verifyEmailLimiter := middleware.RateLimiterMiddleware(10,time.Minute)

	// AUTH ROUTES with rate limits
	r.POST("/signup", signupLimiter, auth.SeekerSignUp)
	r.GET("/verify-email", verifyEmailLimiter, auth.VerifyEmail)
	r.POST("/login", loginLimiter, auth.Login)
	r.POST("/auth/forgot-password", forgotPassLimiter, auth.ForgotPasswordHandler)
	r.POST("/auth/admin-reset-token", auth.SystemInitiatedResetTokenHandler) // (no limiter here)
	r.GET("/reset-password", auth.ResetPasswordPage)                         // (optional)
	r.POST("/reset-password", resetPassLimiter, auth.ResetPasswordHandler)
}
