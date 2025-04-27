package routes

import (
	"RAAS/config"
	"RAAS/middlewares"
	"RAAS/handlers/repo"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"strings"
)

func SetupRoutes(r *gin.Engine, client *mongo.Client, cfg *config.Config) {
	// Set up allowed CORS origins
	origins := strings.Split(cfg.Project.CORSAllowedOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	corsConfig := cors.Config{
		AllowOrigins:  origins,
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}

	// Apply CORS middleware
	r.Use(cors.New(corsConfig))

	// Middleware: Inject MongoDB client into context (adjusted for MongoDB)
	r.Use(middleware.InjectDB(client))

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
	SetupDataEntryRoutes(r, client, cfg)
	SetupFeatureRoutes(r, client, cfg)

	// Reset DB route (update the handler to work with MongoDB)
	r.POST("/api/reset-db", repo.ResetDBHandler)
	r.POST("/api/print-all-collections", repo.PrintAllCollectionsHandler)
}
