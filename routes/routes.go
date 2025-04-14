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

	// Serve static files from the dist folder
	r.Static("/assets", "./public/dist/assets")
	// Serve React app's index.html
	r.GET("/", func(c *gin.Context) {
		c.File("./public/dist/index.html")
	})

	// Catch-all route for client-side routing
	r.NoRoute(func(c *gin.Context) {
		c.File("./public/dist/index.html")
	})

	// Call SetupAuthRoutes, SetupDataEntryRoutes, SetupFeatureRoutes
	SetupAuthRoutes(r, cfg)
	SetupDataEntryRoutes(r, db, cfg)
	SetupFeatureRoutes(r, db, cfg)

	// Reset DB route
	r.POST("/api/reset-db", repo.ResetDBHandler)
}