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

type PersonalInfoRequest struct {
	FirstName       string  `json:"firstName" binding:"required"`
	SecondName      *string `json:"secondName,omitempty"`
	DateOfBirth     string  `json:"dateOfBirth" binding:"required"` // 🛠 Updated to match incoming JSON
	Address         string  `json:"address" binding:"required"`
	LinkedInProfile *string `json:"linkedinProfile,omitempty"`
}

type PersonalInfoResponse struct {
	AuthUserID      uuid.UUID `json:"authUserId"`
	FirstName       string    `json:"firstName"`
	SecondName      *string   `json:"secondName,omitempty"`
	DateOfBirth     string    `json:"dateOfBirth"` // 🛠 Updated to match response key
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

// WorkExperienceRequest struct for incoming work experience data
type WorkExperienceRequest struct {
	JobTitle            string     `json:"jobTitle" binding:"required"`
	CompanyName         string     `json:"companyName" binding:"required"`
	EmploymentType      string     `json:"employmentType" binding:"required"` // e.g., Full-time, Contract
	StartDate           time.Time  `json:"startDate" binding:"required"`       // Format: YYYY-MM-DD
	EndDate             time.Time 	`json:"endDate" binding:"required"`               
	KeyResponsibilities string     `json:"keyResponsibilities" binding:"required"`
}

// WorkExperienceResponse struct for returning work experience data
type WorkExperienceResponse struct {
	ID                  uint       `json:"id"`                          // Dynamically generated ID
	AuthUserID          uuid.UUID  `json:"authUserId"`                  // Associated user ID
	JobTitle            string     `json:"jobTitle"`                    // Job title
	CompanyName         string     `json:"companyName"`                 // Company name
	EmploymentType      string     `json:"employmentType"`              // Employment type (e.g., Full-time, Contract)
	StartDate           time.Time  `json:"startDate"`                   // Start date
	EndDate             time.Time `json:"endDate"`           
	KeyResponsibilities string     `json:"keyResponsibilities"`         // Key responsibilities
}

//EDUCATION
type EducationRequest struct {
	Degree       string     `json:"degree" binding:"required"`
	Institution  string     `json:"institution" binding:"required"`
	FieldOfStudy string     `json:"fieldOfStudy" binding:"required"`
	StartDate    time.Time  `json:"startDate" binding:"required"`
	EndDate      time.Time `json:"endDate,omitempty"`
	Achievements string     `json:"achievements,omitempty"`
}

// EducationResponse represents data sent back to client
type EducationResponse struct {
	ID           uint       `json:"id"`
	AuthUserID   uuid.UUID  `json:"authUserId"`
	Degree       string     `json:"degree"`
	Institution  string     `json:"institution"`
	FieldOfStudy string     `json:"fieldOfStudy"`
	StartDate    time.Time  `json:"startDate"`
	EndDate      time.Time `json:"endDate,omitempty"`
	Achievements string     `json:"achievements,omitempty"`
}

// CertificateRequest represents the input for creating or updating a certificate
type CertificateRequest struct {
	CertificateName   string  `form:"certificateName" json:"certificateName"`
	CertificateNumber *string `form:"certificateNumber" json:"certificateNumber"`
}
// CertificateResponse represents the response sent to the client
type CertificateResponse struct {
	ID                uint      `json:"id"`
	AuthUserID        uuid.UUID `json:"authUserId"`
	CertificateName   string    `json:"certificateName"`
	CertificateFile   string    `json:"certificateFile"`
	CertificateNumber *string   `json:"certificateNumber,omitempty"`
}

//LANGUAGES
type LanguageRequest struct {
	LanguageName     string `json:"language" binding:"required"`
	CertificateFile  string `json:"certificateFile"`
	ProficiencyLevel string `json:"proficiency" binding:"required"`
}

// LanguageResponse is the DTO used for returning language details
type LanguageResponse struct {
	ID               uint      `json:"id"`
	AuthUserID       uuid.UUID `json:"authUserId"`
	LanguageName     string    `json:"language"`
	CertificateFile  string    `json:"certificateFile"`
	ProficiencyLevel string    `json:"proficiency"`
}
