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
}
