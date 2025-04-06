package models

import (
	"gorm.io/gorm"
)
// AUTH MODELS 

type AuthUser struct {
    gorm.Model
    Email              string `gorm:"unique"`
    Password           string // optional for OAuth (can be empty)
    Role               string
    VerificationToken  string
    EmailVerified      bool
    Provider           string `gorm:"default:'local'"` // 'google' or 'local'
}

type Seeker struct {
    gorm.Model
    AuthUserID uint
    AuthUser   AuthUser `gorm:"constraint:OnDelete:CASCADE"`
    FirstName  string
    LastName   string
    Location   string
    // Any other fields specific to Seekers
}


type Admin struct {
    gorm.Model
    AuthUserID uint
    AuthUser   AuthUser `gorm:"constraint:OnDelete:CASCADE"`
    // Any other fields specific to Admins
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