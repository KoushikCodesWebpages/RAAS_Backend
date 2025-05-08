package repository

import (
	"RAAS/internal/models"

	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/google/uuid"

)


// FindSeekerByUserID is a global utility function to find a Seeker by userID in MongoDB
func FindSeekerByUserID(collection *mongo.Collection, userID uuid.UUID) (*models.Seeker, error) {
	var seeker models.Seeker
	filter := bson.M{"auth_user_id": userID}
	err := collection.FindOne(context.Background(), filter).Decode(&seeker)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("seeker not found")
		}
		return nil, err
	}
	return &seeker, nil
}

func IsFieldFilled(personalInfo bson.M) bool {
	// Check if the bson.M map is empty
	return len(personalInfo) > 0
}

func randomSalary() (int, int) {
	// Example random salary range logic, adjust as needed
	minSalary := 25000 // Example minimum salary
	maxSalary := 35000 // Example maximum salary
	return minSalary, maxSalary
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


