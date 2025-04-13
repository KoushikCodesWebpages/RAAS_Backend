package routes

import (
	"RAAS/config"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"RAAS/middlewares"
	"RAAS/handlers/repo"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Middleware: Inject DB into context
	r.Use(middleware.InjectDB(db))

	// Health check route
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>RAAS Backend</title>
			</head>
			<body>
				<h1>RAAS Backend</h1>
				<p>Status: Up and Running</p>
			</body>
			</html>
		`)
	})

	// Call SetupAuthRoutes, SetupDataEntryRoutes, SetupFeatureRoutes
	SetupAuthRoutes(r, cfg)
	SetupDataEntryRoutes(r, db, cfg)
	SetupFeatureRoutes(r, db, cfg)

	// Reset DB route
	r.POST("/api/reset-db", repo.ResetDBHandler)
}
