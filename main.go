package main

import (
	"fmt"
	"log"
	"RAAS/config"
	"RAAS/models"
	"RAAS/routes"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode) 

	// Initialize the configuration
	cfg, err := config.InitConfig()	
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	if cfg.DBType == "sqlite" {
		log.Println("Resetting database...for sqlite")
		// Delete the SQLite database file if it exists
		if _, err := os.Stat(cfg.DBName); err == nil {
			log.Println("SQLite database file exists, deleting it...")
			if err := os.Remove(cfg.DBName); err != nil {
				log.Fatalf("Error deleting SQLite database file: %v", err)
			} else {
				log.Println("SQLite database file deleted successfully.")
			}
		}
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
