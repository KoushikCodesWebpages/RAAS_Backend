package dto

import "github.com/google/uuid"

//JOB TITLE

type JobTitleInput struct {
	PrimaryTitle   string  `json:"primaryTitle"`
	SecondaryTitle *string `json:"secondaryTitle,omitempty"`
	TertiaryTitle  *string `json:"tertiaryTitle,omitempty"`
}


type PersonalInfoRequest struct {
	AuthUserID      uuid.UUID `json:"authUserId"`                 // Set from JWT or session
	FirstName       string    `json:"firstName" binding:"required"`
	SecondName      *string   `json:"secondName,omitempty"`       // Optional
	DateOfBirth     string    `json:"dob" binding:"required"`     // Format: YYYY-MM-DD
	Address         string    `json:"address" binding:"required"`
	LinkedInProfile *string   `json:"linkedinProfile,omitempty"`  // Optional
}

type PersonalInfoResponse struct {
	AuthUserID      uuid.UUID `json:"authUserId"`
	FirstName       string    `json:"firstName"`
	SecondName      *string   `json:"secondName,omitempty"`
	DateOfBirth     string    `json:"dob"`
	Address         string    `json:"address"`
	LinkedInProfile *string   `json:"linkedinProfile,omitempty"`
}
