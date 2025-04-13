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
				<style>
					body {
						font-family: 'Arial', sans-serif;
						background-color: #f4f4f9;
						color: #333;
						margin: 0;
						padding: 0;
						display: flex;
						justify-content: center;
						align-items: center;
						height: 100vh;
					}
					.container {
						text-align: center;
						background-color: #ffffff;
						padding: 30px;
						border-radius: 10px;
						box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
					}
					h1 {
						color: #4CAF50;
						font-size: 3em;
						margin-bottom: 10px;
					}
					p {
						font-size: 1.2em;
						color: #555;
					}
					footer {
						margin-top: 20px;
						font-size: 0.9em;
						color: #777;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<h1>RAAS Backend</h1>
					<p>Status: <span style="color: #28a745;">Up and Running</span></p>
					<footer>
						<p>&copy; 2025 RAAS, Inc. All rights reserved.</p>
					</footer>
				</div>
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
