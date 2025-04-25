package models

import (

	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)


type Job struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	JobID          string             `bson:"jobId" json:"jobId"`
	Title          string             `bson:"title" json:"title"`
	Company        string             `bson:"company" json:"company"`
	Location       string             `bson:"location" json:"location"`
	PostedDate     string             `bson:"postedDate" json:"postedDate"`
	Link           string             `bson:"link" json:"link"`
	Processed      bool               `bson:"processed" json:"processed"`
	Source         string             `bson:"source" json:"source"`

	// Job Description
	JobDescription string             `bson:"jobDescription" json:"jobDescription"`
	JobType        string             `bson:"jobType" json:"jobType"`
	Skills         string             `bson:"skills" json:"skills"`


	JobLink string `bson:"jobLink" json:"jobLink"`

	SelectedCount int `bson:"selectedCount" json:"selectedCount"`
}

func CreateJobIndexes(collection *mongo.Collection) error {
	// Unique index for jobId
	jobIdIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "jobId", Value: 1}}, // Index on jobId, unique
		Options: options.Index().SetUnique(true),
	}

	// Unique index for jobLink
	jobLinkIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "jobLink", Value: 1}}, // Index on jobLink, unique
		Options: options.Index().SetUnique(true),
	}

	// Hashed index for title (for job title-based lookups)
	jobTitleIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: "hashed"}}, // Hashed index on title
		Options: options.Index().SetUnique(false),       // Not unique
	}

	// Hashed index for selectedCount (for fast equality checks)
	selectedCountIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "selectedCount", Value: "hashed"}}, // Hashed index on selectedCount
		Options: options.Index().SetUnique(false),               // Not unique
	}

	// Create indexes
	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		jobIdIndex, jobLinkIndex, jobTitleIndex, selectedCountIndex,
	})
	return err
}

// MatchScore for job seeker match score
type MatchScore struct {
	AuthUserID uuid.UUID `bson:"authUserId" json:"authUserId"`
	JobID      string    `bson:"jobId" json:"jobId"`          
	MatchScore float64   `bson:"matchScore" json:"matchScore"` 
}


func CreateMatchScoreIndexes(collection *mongo.Collection) error {
	// Compound unique index for authUserId and jobId
	matchScoreIndex := mongo.IndexModel{
		Keys:    bson.D{
			{Key: "authUserId", Value: 1}, // Index on authUserId
			{Key: "jobId", Value: 1},      // Index on jobId
		},
		Options: options.Index().SetUnique(true), // Ensuring the combination is unique
	}

	// Create the compound index
	_, err := collection.Indexes().CreateOne(context.Background(), matchScoreIndex)
	return err
}














