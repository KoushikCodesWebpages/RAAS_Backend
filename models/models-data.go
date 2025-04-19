package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AUTH MODELS

type AuthUser struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Email             string    `gorm:"unique;not null"`
	Phone             string    `gorm:"not null"`
	Password          string
	Role              string
	VerificationToken string
	EmailVerified     bool
	Provider          string `gorm:"default:'local'"`
}

// Seeker represents a user who is seeking a job.
type Seeker struct {
    gorm.Model
    AuthUserID uuid.UUID `gorm:"type:char(36);uniqueIndex;not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"authUserId"`

    SubscriptionTier          string         `gorm:"default:'free'" json:"subscriptionTier"`
    DailySelectableJobsCount  int            `gorm:"default:5" json:"dailySelectableJobsCount"`       
    DailyGeneratableCV        int            `gorm:"default:100" json:"dailyGeneratableCv"`             
    DailyGeneratableCoverletter int          `gorm:"default:100" json:"dailyGeneratableCoverletter"`    
    TotalApplications         int            `gorm:"default:0" json:"totalApplications"`              
    
    // Personal Info
    PersonalInfo datatypes.JSON `gorm:"type:json" json:"personalInfo"` // Store JSON for personal info

    // Professional Summary
    ProfessionalSummary datatypes.JSON `gorm:"type:json" json:"professionalSummary"`

    // Work Experience
    WorkExperiences datatypes.JSON `gorm:"type:json" json:"workExperiences"`
	
    // Education
    Educations datatypes.JSON `gorm:"type:json" json:"education"`

    // Certificates
    Certificates datatypes.JSON `gorm:"type:json" json:"certificates"`

    // Languages
    Languages datatypes.JSON `gorm:"type:json" json:"languages"`
    
    // Preferred Job Titles
    PrimaryTitle   string    `gorm:"type:varchar(255);" json:"primaryTitle"`
    SecondaryTitle *string   `gorm:"type:varchar(255);" json:"secondaryTitle"`
    TertiaryTitle  *string   `gorm:"type:varchar(255);" json:"tertiaryTitle"`
}


type Admin struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);unique;not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"authUserId"`
}

// PREFERENCE MODELS

type PersonalInfo struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);unique;not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"authUserId"`
	FirstName       string    `gorm:"type:varchar(100);not null" json:"firstName"`
	SecondName      *string   `gorm:"type:varchar(100)" json:"secondName"`
	DateOfBirth     string    `gorm:"type:date;not null" json:"dateOfBirth"`
	Address         string    `gorm:"type:text;not null" json:"address"`
	LinkedInProfile *string   `gorm:"type:varchar(255)" json:"linkedinProfile"`
}

type ProfessionalSummary struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);unique;not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"authUserId"`
	About        string         `gorm:"type:text;not null" json:"about"`
	Skills       datatypes.JSON `gorm:"type:json;not null" json:"skills"`
	AnnualIncome float64        `gorm:"not null" json:"annualIncome"`
}

type WorkExperience struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"authUserId"`
	JobTitle           string     `gorm:"type:varchar(100);not null" json:"jobTitle"`
	CompanyName        string     `gorm:"type:varchar(100);not null" json:"companyName"`
	EmploymentType       string     `gorm:"type:varchar(50);not null" json:"employerType"`
	StartDate          time.Time  `gorm:"not null" json:"startDate"`
	EndDate            *time.Time `json:"endDate,omitempty"`
	KeyResponsibilities string     `gorm:"type:text" json:"keyResponsibilities"`
}

type Education struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"authUserId"`
	Degree       string     `gorm:"type:varchar(100);not null" json:"degree"`
	Institution  string     `gorm:"type:varchar(150);not null" json:"institution"`
	FieldOfStudy string     `gorm:"type:varchar(100);not null" json:"fieldOfStudy"`
	StartDate    time.Time  `gorm:"not null" json:"startDate"`
	EndDate      *time.Time `json:"endDate,omitempty"`
	Achievements string     `gorm:"type:text" json:"achievements"`
}

type Certificate struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"authUserId"`
	CertificateName   string    `gorm:"type:varchar(255);not null" json:"certificateName"`
	CertificateFile   string    `gorm:"type:text;not null" json:"certificateFile"`
	CertificateNumber string    `gorm:"type:varchar(100)" json:"certificateNumber"`
}

type Language struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"authUserId"`
	LanguageName     string    `gorm:"type:varchar(100);not null" json:"language"`
	CertificateFile  string    `gorm:"type:text" json:"certificateFile"`
	ProficiencyLevel string    `gorm:"type:varchar(20);not null" json:"proficiency"`
}

type PreferredJobTitle struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);unique;not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"authUserId"`
	PrimaryTitle   string    `gorm:"type:varchar(255);not null" json:"primaryTitle"`
	SecondaryTitle *string   `gorm:"type:varchar(255);" json:"secondaryTitle"`
	TertiaryTitle  *string   `gorm:"type:varchar(255);" json:"tertiaryTitle"`
}

// FEATURE MODELS

type JobMatchScore struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	AuthUserID uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:idx_user_job" json:"authUserId"`
	JobID      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_job"` // ✅ Specify length
	Platform   string    `gorm:"type:varchar(50);not null"`                           // Better to limit this too
	Score      float64   `gorm:"not null"`
	MatchedAt  time.Time `gorm:"autoCreateTime"`
}