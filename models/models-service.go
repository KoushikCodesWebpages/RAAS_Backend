package models

import
(
	"github.com/google/uuid"
)

// Job model combining Metadata, Description, and JobLink
type Job struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string `gorm:"unique;type:varchar(191)"`
	Title          string
	Company        string
	Location       string
	PostedDate     string
	Link           string `gorm:"unique;type:varchar(191)"`
	Processed      bool
	Source         string // LinkedIn or Xing (to differentiate the source)

	// Job Description
	JobDescription string
	JobType        string
	Skills         string

	// JobLink for unique reference
	JobLink string `gorm:"unique;type:varchar(191);not null"`
}

	type MatchScore struct {
		SeekerID   uuid.UUID `gorm:"type:char(36);primaryKey"`  // Seeker ID (foreign key reference)
		JobID      string    `gorm:"primaryKey"`                // Job ID (foreign key reference, as string)
		MatchScore float64   `gorm:"type:float"`                // Match score percentage (0 to 100)
	}
























// LinkedIn job metadata
type LinkedInJobMetaData struct {
	ID         string `gorm:"primaryKey;type:varchar(191)"`
	JobID      string `gorm:"unique;type:varchar(191)"`
	Title      string
	Company    string
	Location   string
	PostedDate string
	Link       string `gorm:"unique;type:varchar(191)"`
	Processed  bool

	// Relationships
}

// Xing job metadata
type XingJobMetaData struct {
	ID         string `gorm:"primaryKey;type:varchar(191)"`
	JobID      string `gorm:"unique;type:varchar(191)"`
	Title      string
	Company    string
	Location   string
	PostedDate string
	Link       string `gorm:"unique;type:varchar(191)"`
	Processed  bool
}

// LinkedIn application links
type LinkedInJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null"`
	JobLink string `gorm:"unique;type:varchar(191)"`
}

// Xing application links
type XingJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job_app"`
	JobLink string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job_app"`

}

// LinkedIn job description
type LinkedInJobDescription struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string `gorm:"type:varchar(191);not null;uniqueIndex:idx_linkedin_job"`
	JobLink        string `gorm:"type:varchar(191);not null;uniqueIndex:idx_linkedin_job"`
	JobDescription string
	JobType        string
	Skills         string
}

// Xing job description
type XingJobDescription struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job"`
	JobLink        string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job"`
	JobDescription string
	JobType        string
	Skills         string

}

