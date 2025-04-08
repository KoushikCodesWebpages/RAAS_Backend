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


// JOB DATA MODELS

type LinkedInJobMetaData struct {
	ID         string `gorm:"primaryKey"`  // ID is the main identifier
	JobID      string `gorm:"unique"`      // JobID is unique per job
	Title      string                      // Title of the job
	Company    string                      // Company name
	Location   string                      // Location text
	PostedDate string                      // String representation of the posted date
	Link       string `gorm:"unique"`      // Each job has a unique link
	Processed  bool                        // Whether this job has been processed or not
}

type LinkedInFailedJob struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"` // Auto-increment primary key
	JobID   string                                   // Foreign key to LinkedInJob.ID
	JobLink string `gorm:"unique"`                   // Unique link for tracking

	// Define relationship for foreign key with cascading delete
	LinkedInJob LinkedInJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}

type LinkedInJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string
	JobLink string `gorm:"unique"`

	LinkedInJob LinkedInJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}

type XingJobMetaData struct {
	ID         string `gorm:"primaryKey"`
	JobID      string `gorm:"unique"`
	Title      string
	Company    string
	Location   string
	PostedDate string
	Link       string `gorm:"unique"`
	Processed  bool
}

type XingFailedJob struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string
	JobLink string `gorm:"unique"`

	XingJob XingJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}

type XingJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"uniqueIndex:idx_xing_job_app"` // part of composite unique key
	JobLink string `gorm:"uniqueIndex:idx_xing_job_app"`

	XingJob XingJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}
