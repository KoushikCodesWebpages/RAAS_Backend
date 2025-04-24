package main

import (
    "fmt"
    "log"

    "os"

    "github.com/gin-gonic/gin"

   

   "RAAS/config"
   "RAAS/models"
   "RAAS/routes"


    "gorm.io/gorm"
    "RAAS/workers"
)
func main() {
    gin.SetMode(gin.ReleaseMode)

    err := config.InitConfig()
    if err != nil {
        log.Fatalf("Error initializing config: %v", err)
    }

    db := models.InitDB(config.Cfg)

    // ✅ Start the match score worker properly
    // startMatchScoreWorker(db)

    r := gin.Default()
    routes.SetupRoutes(r, db, config.Cfg)

    port := os.Getenv("PORT")
    if port == "" {
        port = fmt.Sprintf("%d", config.Cfg.ServerPort)
        log.Printf("Starting server on dev port: http://localhost:%s", port)
    } else {
        log.Printf("Starting server on prod port: %s", port)
    }

    err = r.Run(":" + port)
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

// ✅ move this outside main
func startMatchScoreWorker(db *gorm.DB) {
    worker := &workers.MatchScoreWorker{
        DB: db,
    }
    go worker.Run()
}
