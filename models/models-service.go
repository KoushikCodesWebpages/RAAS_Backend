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
	FailedJobs     []LinkedInFailedJob         `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
	ApplicationLinks []LinkedInJobApplicationLink `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
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

	// Relationships
	FailedJobs     []XingFailedJob         `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
	ApplicationLinks []XingJobApplicationLink `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}

// Failed LinkedIn jobs
type LinkedInFailedJob struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null"`
	JobLink string `gorm:"unique;type:varchar(191)"`

	LinkedInJob LinkedInJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}

// Failed Xing jobs
type XingFailedJob struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null"`
	JobLink string `gorm:"unique;type:varchar(191)"`

	XingJob XingJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
}

// LinkedIn application links
type LinkedInJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null"`
	JobLink string `gorm:"unique;type:varchar(191)"`

	LinkedInJob LinkedInJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
	Descriptions []LinkedInJobDescription `gorm:"foreignKey:JobLink;references:JobLink;constraint:OnDelete:CASCADE"`
}

// Xing application links
type XingJobApplicationLink struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	JobID   string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job_app"`
	JobLink string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job_app"`

	XingJob XingJobMetaData `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
	Descriptions []XingJobDescription `gorm:"foreignKey:JobLink;references:JobLink;constraint:OnDelete:CASCADE"`
}

// LinkedIn job description
type LinkedInJobDescription struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string `gorm:"type:varchar(191);not null;uniqueIndex:idx_linkedin_job"`
	JobLink        string `gorm:"type:varchar(191);not null;uniqueIndex:idx_linkedin_job"`
	JobDescription string
	JobType        string
	Skills         string

	JobAppLink LinkedInJobApplicationLink `gorm:"foreignKey:JobLink;references:JobLink;constraint:OnDelete:CASCADE"`
}

// Xing job description
type XingJobDescription struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	JobID          string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job"`
	JobLink        string `gorm:"type:varchar(191);not null;uniqueIndex:idx_xing_job"`
	JobDescription string
	JobType        string
	Skills         string

	JobAppLink XingJobApplicationLink `gorm:"foreignKey:JobLink;references:JobLink;constraint:OnDelete:CASCADE"`
}
