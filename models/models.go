package models

import (
	"gorm.io/gorm"
    "github.com/google/uuid"
)
// AUTH MODELS 

type AuthUser struct {
	gorm.Model
    ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email             string    `gorm:"unique"`
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
	AuthUserID uuid.UUID  `gorm:"type:uuid"`
	AuthUser   AuthUser   `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
	FirstName  string
	LastName   string
	Location   string
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





type PreferredJobTitle struct {
	gorm.Model
	PrimaryTitle   string    `gorm:"type:varchar(255);not null" json:"primaryTitle"`   // Primary job title (cannot be NULL)
	SecondaryTitle *string   `gorm:"type:varchar(255);" json:"secondaryTitle"`         // Secondary job title (nullable)
	TertiaryTitle  *string   `gorm:"type:varchar(255);" json:"tertiaryTitle"`          // Tertiary job title (nullable)
	AuthUserID     uuid.UUID `gorm:"type:uuid;not null" json:"authUserId"`             // Foreign Key to AuthUser
}


// FEATURE MODELS
type LinkedinJobMetadata struct {
	gorm.Model
	JobID      int64  `gorm:"uniqueIndex" json:"jobId"`
	Title      string `gorm:"type:VARCHAR(255)" json:"title"`
	Company    string `gorm:"type:VARCHAR(255)" json:"company"`
	Location   string `gorm:"type:VARCHAR(255)" json:"location"`
	PostedDate string `gorm:"type:DATE" json:"postedDate"`
	Link       string `gorm:"type:VARCHAR(1000);uniqueIndex" json:"link"`
	Processed  bool   `json:"processed"`
}