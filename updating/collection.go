package repository

import (
	"context"
	"RAAS/internal/models"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// Fetch seeker and extract skills
func GetSeekerData(db *mongo.Database, userID string) (models.Seeker, []string, error) {
	var seeker models.Seeker
	err := db.Collection("seekers").FindOne(context.TODO(), bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		return models.Seeker{}, nil, err
	}
	skills := []string{}
	if seeker.ProfessionalSummary != nil {
		skills = extractSkills(seeker.ProfessionalSummary) // Use your existing skill extraction logic
	}
	return seeker, skills, nil
}


// Extract preferred titles from seeker
func CollectPreferredTitles(seeker models.Seeker) []string {
	var titles []string
	if seeker.PrimaryTitle != "" {
		titles = append(titles, seeker.PrimaryTitle)
	}
	if seeker.SecondaryTitle != nil && *seeker.SecondaryTitle != "" {
		titles = append(titles, *seeker.SecondaryTitle)
	}
	if seeker.TertiaryTitle != nil && *seeker.TertiaryTitle != "" {
		titles = append(titles, *seeker.TertiaryTitle)
	}
	return titles
}

// Fetch job by job ID
func GetJobByID(db *mongo.Database, jobID string) (models.Job, error) {
	var job models.Job
	err := db.Collection("jobs").FindOne(context.TODO(), bson.M{"job_id": jobID}).Decode(&job)
	if err != nil {
		return models.Job{}, err
	}
	return job, nil
}

// Extract skills safely
func extractSkills(professionalSummary bson.M) []string {
	if val, ok := professionalSummary["skills"].(primitive.A); ok {
		var skills []string
		for _, skill := range val {
			if str, ok := skill.(string); ok {
				skills = append(skills, str)
			}
		}
		return skills
	}
	return nil
}

// Helper function to fetch saved job IDs
func FetchSavedJobIDs(c *gin.Context, col *mongo.Collection, userID string) ([]string, error) {
	var jobIDs []string
	cursor, err := col.Find(c, bson.M{"auth_user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var saved models.SavedJob
		if err := cursor.Decode(&saved); err == nil {
			jobIDs = append(jobIDs, saved.JobID)
		}
	}
	return jobIDs, nil
}

