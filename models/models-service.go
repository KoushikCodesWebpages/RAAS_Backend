package models

import (
	// "gorm.io/gorm"
    // "github.com/google/uuid"
	// "time"
	// "gorm.io/datatypes"

	
)

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





type LinkedInFailedJob struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"` // Auto-increment primary key
	JobID   string                                   // Foreign key to LinkedInJob.ID
	JobLink string `gorm:"unique"`                   // Unique link for tracking

	// Define relationship for foreign key with cascading delete
	LinkedInJob LinkedInJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}
type XingFailedJob struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string
	JobLink string `gorm:"unique"`

	XingJob XingJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}



type LinkedInJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string
	JobLink string `gorm:"unique"`

	LinkedInJob LinkedInJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}
type XingJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"uniqueIndex:idx_xing_job_app"` // part of composite unique key
	JobLink string `gorm:"uniqueIndex:idx_xing_job_app"`

	XingJob XingJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}



type LinkedInJobDescription struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string
	JobLink        string `gorm:"unique"`
	JobDescription string

	LinkedInJob LinkedInJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}
type XingJobDescription struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string
	JobLink        string `gorm:"unique"`
	JobDescription string

	XingJob XingJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}



