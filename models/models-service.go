package models

import (
	// "gorm.io/gorm"
	// "time"
)

// JOB DATA MODELS

// Job metadata from LinkedIn
type LinkedInJobMetaData struct {
	ID         string `gorm:"primaryKey"`  // Unique primary key
	JobID      string `gorm:"unique;type:varchar(191)"`      // Unique job identifier from LinkedIn, with fixed length
	Title      string                      // Job title
	Company    string                      // Company name
	Location   string                      // Job location
	PostedDate string                      // Posted date (string format)
	Link       string `gorm:"unique;type:varchar(191)"`      // Unique job link, with fixed length
	Processed  bool                        // Whether the job has been processed
}

// Job metadata from Xing
type XingJobMetaData struct {
	ID         string `gorm:"primaryKey"`
	JobID      string `gorm:"unique;type:varchar(191)"`      // Unique job identifier from Xing, with fixed length
	Title      string
	Company    string
	Location   string
	PostedDate string
	Link       string `gorm:"unique;type:varchar(191)"`      // Unique job link, with fixed length
	Processed  bool
}

// Failed job records for LinkedIn
type LinkedInFailedJob struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null"`     // Fixed length for job ID
	JobLink string `gorm:"unique;type:varchar(191)"`      // Unique job link

}

// Failed job records for Xing
type XingFailedJob struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null"`
	JobLink string `gorm:"unique;type:varchar(191)"`

}

// Application links for LinkedIn jobs
type LinkedInJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null"`
	JobLink string `gorm:"unique;type:varchar(191)"`

	LinkedInJob LinkedInJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}

// Application links for Xing jobs
type XingJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"uniqueIndex:idx_xing_job_app;type:varchar(191)"`
	JobLink string `gorm:"uniqueIndex:idx_xing_job_app;type:varchar(191)"`

}

type LinkedInJobDescription struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string `gorm:"type:varchar(191);not null;uniqueIndex:idx_linkedin_job"`
	JobLink        string `gorm:"type:varchar(191);not null;uniqueIndex:idx_linkedin_job"`
	JobDescription string
	JobType        string
	Skills         string
}

type XingJobDescription struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job"`
	JobLink        string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job"`
	JobDescription string
	JobType        string
	Skills         string
}
