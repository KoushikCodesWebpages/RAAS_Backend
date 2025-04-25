package dto

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// =======================
// JOB TITLE
// =======================

type JobTitleInput struct {
	PrimaryTitle   string  `json:"primaryTitle" bson:"primaryTitle"`
	SecondaryTitle *string `json:"secondaryTitle,omitempty" bson:"secondaryTitle,omitempty"`
	TertiaryTitle  *string `json:"tertiaryTitle,omitempty" bson:"tertiaryTitle,omitempty"`
}

type JobTitleResponse struct {
	AuthUserID     uuid.UUID `json:"authUserId" bson:"authUserId"`
	PrimaryTitle   string    `json:"primaryTitle" bson:"primaryTitle"`
	SecondaryTitle *string   `json:"secondaryTitle,omitempty" bson:"secondaryTitle,omitempty"`
	TertiaryTitle  *string   `json:"tertiaryTitle,omitempty" bson:"tertiaryTitle,omitempty"`
}

// =======================
// PERSONAL INFO
// =======================

type PersonalInfoRequest struct {
	FirstName       string  `json:"firstName" binding:"required" bson:"firstName"`
	SecondName      *string `json:"secondName,omitempty" bson:"secondName,omitempty"`
	DateOfBirth     string  `json:"dateOfBirth" binding:"required" bson:"dateOfBirth"`
	Address         string  `json:"address" binding:"required" bson:"address"`
	LinkedInProfile *string `json:"linkedinProfile,omitempty" bson:"linkedinProfile,omitempty"`
}

type PersonalInfoResponse struct {
	AuthUserID      uuid.UUID `json:"authUserId" bson:"authUserId"`
	FirstName       string    `json:"firstName" bson:"firstName"`
	SecondName      *string   `json:"secondName,omitempty" bson:"secondName,omitempty"`
	DateOfBirth     string    `json:"dateOfBirth" bson:"dateOfBirth"`
	Address         string    `json:"address" bson:"address"`
	LinkedInProfile *string   `json:"linkedinProfile,omitempty" bson:"linkedinProfile,omitempty"`
}

// =======================
// PROFESSIONAL SUMMARY
// =======================

type ProfessionalSummaryRequest struct {
	About        string   `json:"about" binding:"required" bson:"about"`
	Skills       []string `json:"skills" binding:"required" bson:"skills"`
	AnnualIncome float64  `json:"annualIncome" binding:"required" bson:"annualIncome"`
}

type ProfessionalSummaryResponse struct {
	AuthUserID   uuid.UUID `json:"authUserId" bson:"authUserId"`
	About        string    `json:"about" bson:"about"`
	Skills       []string  `json:"skills" bson:"skills"`
	AnnualIncome float64   `json:"annualIncome" bson:"annualIncome"`
}

// =======================
// WORK EXPERIENCE
// =======================

type WorkExperienceRequest struct {
	JobTitle            string    `json:"jobTitle" binding:"required" bson:"jobTitle"`
	CompanyName         string    `json:"companyName" binding:"required" bson:"companyName"`
	EmploymentType      string    `json:"employmentType" binding:"required" bson:"employmentType"`
	StartDate           time.Time `json:"startDate" binding:"required" bson:"startDate"` // Format: YYYY-MM-DD
	EndDate             *time.Time `json:"endDate,omitempty" bson:"endDate,omitempty"`    // Optional for ongoing jobs
	KeyResponsibilities string    `json:"keyResponsibilities" binding:"required" bson:"keyResponsibilities"`
}

type WorkExperienceResponse struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID          uuid.UUID          `json:"authUserId" bson:"authUserId"`
	JobTitle            string             `json:"jobTitle" bson:"jobTitle"`
	CompanyName         string             `json:"companyName" bson:"companyName"`
	EmploymentType      string             `json:"employmentType" bson:"employmentType"`
	StartDate           time.Time          `json:"startDate" bson:"startDate"`
	EndDate             *time.Time         `json:"endDate,omitempty" bson:"endDate,omitempty"`
	KeyResponsibilities string             `json:"keyResponsibilities" bson:"keyResponsibilities"`
}

// =======================
// EDUCATION
// =======================

type EducationRequest struct {
	Degree       string     `json:"degree" binding:"required" bson:"degree"`
	Institution  string     `json:"institution" binding:"required" bson:"institution"`
	FieldOfStudy string     `json:"fieldOfStudy" binding:"required" bson:"fieldOfStudy"`
	StartDate    time.Time  `json:"startDate" binding:"required" bson:"startDate"`
	EndDate      *time.Time `json:"endDate,omitempty" bson:"endDate,omitempty"`
	Achievements string     `json:"achievements,omitempty" bson:"achievements,omitempty"`
}

type EducationResponse struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID   uuid.UUID          `json:"authUserId" bson:"authUserId"`
	Degree       string             `json:"degree" bson:"degree"`
	Institution  string             `json:"institution" bson:"institution"`
	FieldOfStudy string             `json:"fieldOfStudy" bson:"fieldOfStudy"`
	StartDate    time.Time          `json:"startDate" bson:"startDate"`
	EndDate      *time.Time         `json:"endDate,omitempty" bson:"endDate,omitempty"`
	Achievements string             `json:"achievements,omitempty" bson:"achievements,omitempty"`
}

// =======================
// CERTIFICATE
// =======================

type CertificateRequest struct {
	CertificateName   string  `form:"certificateName" json:"certificateName" bson:"certificateName"`
	CertificateNumber *string `form:"certificateNumber" json:"certificateNumber,omitempty" bson:"certificateNumber,omitempty"`
}

type CertificateResponse struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID        uuid.UUID          `json:"authUserId" bson:"authUserId"`
	CertificateName   string             `json:"certificateName" bson:"certificateName"`
	CertificateFile   string             `json:"certificateFile" bson:"certificateFile"`
	CertificateNumber *string            `json:"certificateNumber,omitempty" bson:"certificateNumber,omitempty"`
}

// =======================
// LANGUAGES
// =======================

type LanguageRequest struct {
	LanguageName     string `json:"language" binding:"required" bson:"language"`
	CertificateFile  string `json:"certificateFile" bson:"certificateFile"`
	ProficiencyLevel string `json:"proficiency" binding:"required" bson:"proficiency"`
}

type LanguageResponse struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID       uuid.UUID          `json:"authUserId" bson:"authUserId"`
	LanguageName     string             `json:"language" bson:"language"`
	CertificateFile  string             `json:"certificateFile" bson:"certificateFile"`
	ProficiencyLevel string             `json:"proficiency" bson:"proficiency"`
}
