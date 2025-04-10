package main

import (
	"fmt"
	"log"
	"RAAS/config"
	"RAAS/models"
	"RAAS/routes"
	"os"
	"github.com/gin-gonic/gin"

)

func main() {
	gin.SetMode(gin.ReleaseMode) 

	// Initialize the configuration
	cfg, err := config.InitConfig()	
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// Initialize the database
	db := models.InitDB(cfg)
	

	// Create a new Gin router
	r := gin.Default()

	// Register all routes
	routes.SetupRoutes(r, db, cfg)

	// Start the server
	port := os.Getenv("PORT")
	

	if port == "" {
		port = fmt.Sprintf("%d", cfg.ServerPort)
		log.Printf("Starting server on dev port: %s", port)
	} else {
		log.Printf("Starting server on prod port: %s", port)
	}
	err = r.Run(":" + port)
	


	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}


}
