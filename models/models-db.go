package models

import (
	"context"
	"log"
	"RAAS/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

// InitDB initializes the MongoDB connection and returns the client and database instances
func InitDB(cfg *config.Config) (*mongo.Client, *mongo.Database) {
	clientOptions := options.Client().ApplyURI(cfg.Cloud.MongoDBUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("❌ Error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("❌ Error pinging MongoDB: %v", err)
	}

	MongoDB = client.Database(cfg.Cloud.MongoDBName)
	log.Println("✅ MongoDB connection established")
	PrintAllCollections()

	// Optionally reset collections
	resetCollections()


	return client, MongoDB
}

// Reset collections if necessary
func resetCollections() {
	collections := []string{
		"auth_users",
		"seekers",
		"admins",
		"match_scores",
		"user_entry_timelines",
		"selected_job_applications",
		"cover_letters",
		"cv",
	}

	for _, col := range collections {
		err := MongoDB.Collection(col).Drop(context.TODO())
		if err != nil {
			log.Printf("⚠️ Error resetting collection %s: %v", col, err)
		} else {
			log.Printf("✅ Collection %s reset", col)
		}
	}
}

// Print all collections
func PrintAllCollections() {
	collections, err := MongoDB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("❌ Error fetching collection names: %v", err)
	}

	log.Println("📦 Collections in the database:")
	for _, col := range collections {
		log.Println(" -", col)
	}
}
