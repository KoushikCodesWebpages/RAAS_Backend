package main

import (
	"fmt"
	"log"
	"RAAS/config"
	"RAAS/models"
	"RAAS/routes"
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
	err = r.Run(fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
