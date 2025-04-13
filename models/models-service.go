package models

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
