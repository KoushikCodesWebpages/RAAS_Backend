package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"

	"github.com/gin-gonic/gin"
	"RAAS/config"
	"RAAS/routes"
	// "RAAS/workers"
	"RAAS/models" // Import the models package to use InitDB

	// "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Initialize the configuration
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// Initialize MongoDB client and database using models.InitDB
	client, _ := models.InitDB(config.Cfg) // Get both client and database

	// ✅ Start the match score worker properly
	// startMatchScoreWorker(client)

	// Set up the Gin router
	r := gin.Default()
	routes.SetupRoutes(r, client, config.Cfg) // Pass client (mongo.Client) to SetupRoutes

	// Get the server port from config or environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", config.Cfg.Server.ServerPort)
		log.Printf("Starting server on dev port: http://localhost:%s", port)
	} else {
		log.Printf("Starting server on prod port: %s", port)
	}

	// Run the Gin server
	go func() {
		err = r.Run(":" + port)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown handling for server and MongoDB connection
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownSignal

	log.Println("Shutting down server...")

	// Close MongoDB client gracefully
	err = client.Disconnect(context.TODO()) // Disconnect the client
	if err != nil {
		log.Fatalf("Error disconnecting MongoDB: %v", err)
	}

	log.Println("MongoDB connection closed gracefully")
}

// // startMatchScoreWorker starts the match score worker as a goroutine
// func startMatchScoreWorker(client *mongo.Client) {
// 	worker := &workers.MatchScoreWorker{
// 		Client: client,
// 	}
// 	go worker.Run()
// }
