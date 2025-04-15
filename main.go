package main

import (
    "fmt"
    "log"

    "os"

    "github.com/gin-gonic/gin"

   

   "RAAS/config"
   "RAAS/models"
   "RAAS/routes"


    // "gorm.io/gorm"
    // "RAAS/workers"
)

func main() {
    gin.SetMode(gin.ReleaseMode)

  

    // Initialize the configuration
    err := config.InitConfig()
    if err != nil {
        log.Fatalf("Error initializing config: %v", err)
    }

    // Initialize the database
    db := models.InitDB(config.Cfg)

    // Start the match score worker in the background
    //go startMatchScoreWorker(db)

    // Create a new Gin router
    r := gin.Default()

    // Register all routes
    routes.SetupRoutes(r, db, config.Cfg)

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = fmt.Sprintf("%d", config.Cfg.ServerPort)
        log.Printf("Starting server on dev port: %s", port)
    } else {
        log.Printf("Starting server on prod port: %s", port)
    }
    err = r.Run(":" + port)
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

// // Start the match score worker
// func startMatchScoreWorker(db *gorm.DB) {
//  worker := &workers.MatchScoreWorker{
//      DB: db,
//  }

//  // Start the worker's Run method in a separate goroutine
//  go worker.Run()
// }