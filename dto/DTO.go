package dto

import (
    "RAAS/models"
	"github.com/google/uuid"
)

type SeekerSignUpInput struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
	Number string `json:"number" binding:"required,len=10"`
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
}

// SeekerResponse represents the response structure for Seeker details
type SeekerResponse struct {
	ID         uint            `json:"id"`       // Use uuid.UUID for Seeker ID
	AuthUserID uuid.UUID       `json:"authUserId"`  // Use uuid.UUID for AuthUserID
	AuthUser   AuthUserMinimal `json:"authUser"`
	//SubscriptionTier string
}

// NewSeekerResponse creates a new SeekerResponse from a Seeker model
func NewSeekerResponse(seeker models.Seeker) SeekerResponse {
	return SeekerResponse{
		ID:         seeker.ID,  // ID as uuid.UUID
		AuthUserID: seeker.AuthUserID,  // AuthUserID as uuid.UUID
		AuthUser: AuthUserMinimal{
			Email:         seeker.AuthUser.Email,
			EmailVerified: seeker.AuthUser.EmailVerified,
			Provider:      seeker.AuthUser.Provider,
		},
	}
}
