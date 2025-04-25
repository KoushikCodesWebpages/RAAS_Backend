package dto

import (
    "RAAS/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type SeekerSignUpInput struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
	Number    string `json:"number" binding:"required,len=10"`
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// AuthUserMinimal represents minimal user details for response
type AuthUserMinimal struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	Provider      string `json:"provider"`
	Number        string `json:"number" binding:"required,len=10"`
}

// SeekerResponse represents the response structure for Seeker details
type SeekerResponse struct {
	ID               primitive.ObjectID `json:"id"`        // Use primitive.ObjectID for Seeker ID in MongoDB
	AuthUserID       primitive.ObjectID `json:"authUserId"` // Use primitive.ObjectID for AuthUserID in MongoDB
	AuthUser         AuthUserMinimal    `json:"authUser"`
	SubscriptionTier string              `json:"subscriptionTier"`
}

func SeekerProfileResponse(seeker models.Seeker) SeekerResponse {
	// Convert uuid.UUID to primitive.ObjectID
	authUserID, _ := primitive.ObjectIDFromHex(seeker.AuthUserID.String())

	return SeekerResponse{
		ID:               seeker.ID,             // ID as primitive.ObjectID
		AuthUserID:       authUserID,            // AuthUserID as primitive.ObjectID
		SubscriptionTier: seeker.SubscriptionTier,
	}
}