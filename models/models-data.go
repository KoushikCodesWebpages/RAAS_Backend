package models

import (
	"gorm.io/gorm"
    "github.com/google/uuid"
	"time"
	"gorm.io/datatypes"

	
)
// AUTH MODELS 

type AuthUser struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email             string    `gorm:"unique;not null"`
	Phone             string    `gorm:"not null"` // Add this line
	Password          string    // optional for OAuth (can be empty)
	Role              string
	VerificationToken string
	EmailVerified     bool
	Provider          string `gorm:"default:'local'"` // 'google' or 'local'
}


// Seeker struct that references AuthUser by UUID
type Seeker struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey"`
	AuthUserID uuid.UUID  `gorm:"type:uuid;not null"`
	AuthUser   AuthUser   `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
	SubscriptionTier string `gorm:"default:'free'" json:"subscriptionTier"` // 'free', 'premium'

	// Any other fields specific to Seekers
}


type Admin struct {
    gorm.Model
    AuthUserID uint
    AuthUser   AuthUser `gorm:"constraint:OnDelete:CASCADE"`
    // Any other fields specific to Admins
}

// PREFERENCE MODELS

type PersonalInfo struct {
	gorm.Model
	AuthUserID      uuid.UUID `gorm:"type:uuid;unique;not null" json:"authUserId"` // 1:1 relationship
	FirstName       string    `gorm:"type:varchar(100);not null" json:"firstName"`
	SecondName      *string   `gorm:"type:varchar(100)" json:"secondName"`
	DateOfBirth     string    `gorm:"type:date;not null" json:"dob"`
	Address         string    `gorm:"type:text;not null" json:"address"`
	LinkedInProfile *string   `gorm:"type:varchar(255)" json:"linkedinProfile"`

	AuthUser AuthUser `gorm:"constraint:OnDelete:CASCADE;"` // FK to AuthUser
}

type ProfessionalSummary struct {
	gorm.Model
	AuthUserID   uuid.UUID `gorm:"type:uuid;unique;not null" json:"authUserId"` // 1:1 with AuthUser

	About        string         `gorm:"type:text;not null" json:"about"`
	Skills       datatypes.JSON `gorm:"type:jsonb;not null" json:"skills"`     // Store as JSON array
	AnnualIncome float64        `gorm:"not null" json:"annualIncome"`

	AuthUser AuthUser `gorm:"constraint:OnDelete:CASCADE;"` // FK to AuthUser
}





type PreferredJobTitle struct {
	gorm.Model
	AuthUserID     uuid.UUID `gorm:"type:uuid;unique;not null" json:"authUserId"`
	PrimaryTitle   string    `gorm:"type:varchar(255);not null" json:"primaryTitle"`   // Primary job title (cannot be NULL)
	SecondaryTitle *string   `gorm:"type:varchar(255);" json:"secondaryTitle"`         // Secondary job title (nullable)
	TertiaryTitle  *string   `gorm:"type:varchar(255);" json:"tertiaryTitle"`          // Tertiary job title (nullable)
	
	AuthUser       AuthUser  `gorm:"constraint:OnDelete:CASCADE;"`       // Foreign Key to AuthUser
}


//FEATURE MODELS

type JobMatchScore struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	AuthUserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_job" json:"authUserId"`
	JobID      string    `gorm:"not null;uniqueIndex:idx_user_job"` // composite with AuthUserID
	Platform   string    `gorm:"not null"`                           // "linkedin" or "xing"
	Score      float64   `gorm:"not null"`
	MatchedAt  time.Time `gorm:"autoCreateTime"`
	AuthUser   AuthUser  `gorm:"constraint:OnDelete:CASCADE;"`
}


