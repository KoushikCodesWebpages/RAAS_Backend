package handlers

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/google/uuid"
	"RAAS/models"
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

