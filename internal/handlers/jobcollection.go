package handlers

import (
	"context"
	"RAAS/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
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

// Fetch job by job ID
func GetJobByID(db *mongo.Database, jobID string) (models.Job, error) {
	var job models.Job
	err := db.Collection("jobs").FindOne(context.TODO(), bson.M{"job_id": jobID}).Decode(&job)
	if err != nil {
		return models.Job{}, err
	}
	return job, nil
}

// Generate expected salary range
func GenerateSalaryRange() models.SalaryRange {
	min, max := randomSalary() // Use your existing function
	return models.SalaryRange{Min: min, Max: max}
}

func dereferenceString(str *string) string {
	if str != nil {
		return *str
	}
	return "" // Return an empty string if the pointer is nil
}


// Helper function to get optional fields
func getOptionalField(info bson.M, field string) *string {
	if val, ok := info[field]; ok && val != nil {
		v := val.(string)
		return &v
	}
	return nil
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

func randomSalary() (int, int) {
	// Example random salary range logic, adjust as needed
	minSalary := 25000 // Example minimum salary
	maxSalary := 35000 // Example maximum salary
	return minSalary, maxSalary
}