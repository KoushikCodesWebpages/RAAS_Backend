package models

import (
	"gorm.io/gorm"
)

// DB is the global database variable

// InitDB initializes the database connection and sets up the models

// AutoMigrate will automatically migrate all models to the database

// ResetDB will truncate tables, reset auto increment, and delete data, but keep table structure

// AUTH MODELS 



type AuthUser struct {
    gorm.Model
    Email    string `gorm:"unique"`
    Password string
    Role     string
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