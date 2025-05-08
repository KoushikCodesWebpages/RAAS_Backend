package dto

import (

	
	"RAAS/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

)


// =======================
// PERSONAL INFO
// =======================

type PersonalInfoRequest struct {
	FirstName       string  `json:"first_name" binding:"required" bson:"first_name"`
	SecondName      *string `json:"second_name,omitempty" bson:"second_name,omitempty"`
	DateOfBirth     string  `json:"date_of_birth" binding:"required" bson:"date_of_birth"`
	Address         string  `json:"address" binding:"required" bson:"address"`
	LinkedInProfile *string `json:"linkedin_profile,omitempty" bson:"linkedin_profile,omitempty"`
}

type PersonalInfoResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID      string              `json:"auth_user_id" bson:"auth_user_id"`
	FirstName       string              `json:"first_name" bson:"first_name"`
	SecondName      *string             `json:"second_name,omitempty" bson:"second_name,omitempty"`
	DateOfBirth     string              `json:"date_of_birth" bson:"date_of_birth"`
	Address         string              `json:"address" bson:"address"`
	LinkedInProfile *string             `json:"linkedin_profile,omitempty" bson:"linkedin_profile,omitempty"`
}

// =======================
// PROFESSIONAL SUMMARY
// =======================

type ProfessionalSummaryRequest struct {
	About        string   `json:"about" binding:"required" bson:"about"`
	Skills       []string `json:"skills" binding:"required" bson:"skills"`
	AnnualIncome float64  `json:"annual_income" binding:"required" bson:"annual_income"`
}

type ProfessionalSummaryResponse struct {
	AuthUserID   string   `json:"auth_user_id" bson:"auth_user_id"`
	About        string   `json:"about" bson:"about"`
	Skills       []string `json:"skills" bson:"skills"`
	AnnualIncome float64  `json:"annual_income" bson:"annual_income"`
}

// =======================
// WORK EXPERIENCE
// =======================

// Request payload (for creating or updating)
type WorkExperienceRequest struct {
	JobTitle            string          `json:"job_title" binding:"required" bson:"job_title"`
	CompanyName         string          `json:"company_name" binding:"required" bson:"company_name"`
	EmploymentType      string          `json:"employment_type" binding:"required" bson:"employment_type"`
	StartDate           utils.DateOnly  `json:"start_date" binding:"required" bson:"start_date"`
	EndDate             *utils.DateOnly `json:"end_date,omitempty" bson:"end_date,omitempty"`
	KeyResponsibilities string          `json:"key_responsibilities" binding:"required" bson:"key_responsibilities"`
}

// Response payload
type WorkExperienceResponse struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID          string              `json:"auth_user_id" bson:"auth_user_id"`
	JobTitle            string              `json:"job_title" bson:"job_title"`
	CompanyName         string              `json:"company_name" bson:"company_name"`
	EmploymentType      string              `json:"employment_type" bson:"employment_type"`
	StartDate           utils.DateOnly      `json:"start_date" bson:"start_date"`
	EndDate             *utils.DateOnly     `json:"end_date,omitempty" bson:"end_date,omitempty"`
	KeyResponsibilities string              `json:"key_responsibilities" bson:"key_responsibilities"`
}

// =======================
// EDUCATION
// =======================

type EducationRequest struct {
	Degree       string     `json:"degree" binding:"required" bson:"degree"`
	Institution  string     `json:"institution" binding:"required" bson:"institution"`
	FieldOfStudy string     `json:"field_of_study" binding:"required" bson:"field_of_study"`
	StartDate    utils.DateOnly  `json:"start_date" binding:"required" bson:"start_date"`
	EndDate      *utils.DateOnly `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Achievements string     `json:"achievements,omitempty" bson:"achievements,omitempty"`
}

type EducationResponse struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID   string              `json:"auth_user_id" bson:"auth_user_id"`
	Degree       string              `json:"degree" bson:"degree"`
	Institution  string              `json:"institution" bson:"institution"`
	FieldOfStudy string              `json:"field_of_study" bson:"field_of_study"`
	StartDate    utils.DateOnly      `json:"start_date" bson:"start_date"`
	EndDate      *utils.DateOnly     `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Achievements string              `json:"achievements,omitempty" bson:"achievements,omitempty"`
}

// =======================
// CERTIFICATE
// =======================

type CertificateRequest struct {
	CertificateName   string  `form:"certificate_name" json:"certificate_name" bson:"certificate_name"`
	CertificateNumber *string `form:"certificate_number" json:"certificate_number,omitempty" bson:"certificate_number,omitempty"`
}

type CertificateResponse struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID        string              `json:"auth_user_id" bson:"auth_user_id"`
	CertificateName   string              `json:"certificate_name" bson:"certificate_name"`
	CertificateFile   string              `json:"certificate_file" bson:"certificate_file"`
	CertificateNumber *string             `json:"certificate_number,omitempty" bson:"certificate_number,omitempty"`
}

// =======================
// LANGUAGES
// =======================

type LanguageRequest struct {
	LanguageName     string `json:"language" binding:"required" bson:"language"`
	CertificateFile  string `json:"certificate_file" bson:"certificate_file"`
	ProficiencyLevel string `json:"proficiency" binding:"required" bson:"proficiency"`
}

type LanguageResponse struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthUserID       string              `json:"auth_user_id" bson:"auth_user_id"`
	LanguageName     string              `json:"language" bson:"language"`
	CertificateFile  string              `json:"certificate_file" bson:"certificate_file"`
	ProficiencyLevel string              `json:"proficiency" bson:"proficiency"`
}

// =======================
// JOB TITLE
// =======================

type JobTitleInput struct {
	PrimaryTitle   string  `json:"primary_title" bson:"primary_title"`
	SecondaryTitle *string `json:"secondary_title,omitempty" bson:"secondary_title,omitempty"`
	TertiaryTitle  *string `json:"tertiary_title,omitempty" bson:"tertiary_title,omitempty"`
}

type JobTitleResponse struct {
	AuthUserID     string  `json:"auth_user_id" bson:"auth_user_id"`
	PrimaryTitle   string  `json:"primary_title" bson:"primary_title"`
	SecondaryTitle *string `json:"secondary_title,omitempty" bson:"secondary_title,omitempty"`
	TertiaryTitle  *string `json:"tertiary_title,omitempty" bson:"tertiary_title,omitempty"`
}
