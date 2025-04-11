package dto

import (
	"github.com/google/uuid"
	"time"
)

//JOB TITLE

type JobTitleInput struct {
	PrimaryTitle   string  `json:"primaryTitle"`
	SecondaryTitle *string `json:"secondaryTitle,omitempty"`
	TertiaryTitle  *string `json:"tertiaryTitle,omitempty"`
}

//PERSONAL INFO
type PersonalInfoRequest struct {                // Set from JWT or session
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

//PROFESSNAL SUMMARY
type ProfessionalSummaryRequest struct {
	About        string   `json:"about" binding:"required"`
	Skills       []string `json:"skills" binding:"required"`
	AnnualIncome float64  `json:"annualIncome" binding:"required"`
}

type ProfessionalSummaryResponse struct {
	AuthUserID   uuid.UUID `json:"authUserId"`
	About        string    `json:"about"`
	Skills       []string  `json:"skills"`
	AnnualIncome float64   `json:"annualIncome"`
}

//WORKEXPERIENCE

type WorkExperienceRequest struct {
	JobTitle            string     `json:"jobTitle" binding:"required"`
	CompanyName         string     `json:"companyName" binding:"required"`
	EmployerType        string     `json:"employerType" binding:"required"` // e.g., Full-time, Contract
	StartDate           time.Time  `json:"startDate" binding:"required"`    // Format: YYYY-MM-DD
	EndDate             *time.Time `json:"endDate,omitempty"`               // Nullable
	KeyResponsibilities string     `json:"keyResponsibilities" binding:"required"`
}

type WorkExperienceResponse struct {
	ID                 uint       `json:"id"`
	AuthUserID         uuid.UUID  `json:"authUserId"`
	JobTitle           string     `json:"jobTitle"`
	CompanyName        string     `json:"companyName"`
	EmployerType       string     `json:"employerType"`
	StartDate          time.Time  `json:"startDate"`
	EndDate            *time.Time `json:"endDate,omitempty"`
	KeyResponsibilities string    `json:"keyResponsibilities"`

}