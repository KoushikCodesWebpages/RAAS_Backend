package dto

import (
    // "RAAS/models"
	"github.com/google/uuid"
	//"time"
    "gorm.io/gorm"
)


//JOB RETRIEVAL


type SalaryRange struct {
    Min int `json:"min"`
    Max int `json:"max"`
}
type JobDTO struct {
    Source         string       `json:"source"`                   // "linkedin" or "xing"
    ID             uint       `json:"id"`                       // UUID or unique DB ID
    JobID          string       `json:"job_id"`                   // Platform-specific Job ID
    Title          string       `json:"title"`
    Company        string       `json:"company"`
    Location       string       `json:"location"`
    PostedDate     string       `json:"posted_date"`
    Processed      bool         `json:"processed"`
    JobType        string       `json:"job_type"`                 // e.g., Full-time, Part-time
    Skills         string       `json:"skills"`                   // Comma-separated required skills
    UserSkills     []string     `json:"userSkills"`               // List of user skills used in matching
    ExpectedSalary SalaryRange  `json:"expected_salary"`          // Expected salary range
    MatchScore     float64      `json:"match_score"`              // Match score from 0 to 100
    Description    string       `json:"description"`              // Job description text
}


// JobFilterDTO represents the filter data for job retrieval.
type JobFilterDTO struct {
	Title string `form:"title"` // Query param: /jobs/linkedin?title=developer
}

// MatchScoreResponse represents the response containing the match score for a job.
type MatchScoreResponse struct {
	SeekerID   uuid.UUID `json:"seeker_id"`
	JobID      string `json:"job_id"`
	MatchScore float64    `json:"match_score"`
}

type SelectJobRequest struct {
    Source         string      `json:"source" binding:"required"`
    ID             string      `json:"id" binding:"required"`          // Unique DB or UUID
    JobID          string      `json:"job_id" binding:"required"`
    Title          string      `json:"title" binding:"required"`
    Company        string      `json:"company"`
    Location       string      `json:"location"`
    PostedDate     string      `json:"posted_date"`
    Processed      bool        `json:"processed"`
    JobType        string      `json:"job_type"`
    Skills         string      `json:"skills"`
    UserSkills     []string    `json:"userSkills"`
    ExpectedSalary SalaryRange `json:"expected_salary"`
    MatchScore     float64     `json:"match_score"`
    Description    string      `json:"description"`
}

type SelectedJobApplication struct {
    gorm.Model
    AuthUserID            uuid.UUID `gorm:"type:char(36);not null" json:"auth_user_id"`  // Foreign key to Seeker
    Source                string    `gorm:"type:varchar(20);not null" json:"source"`     // "linkedin" or "xing"
    JobID                 string    `gorm:"type:varchar(100);not null" json:"job_id"`    // Platform-specific job ID
    Title                 string    `gorm:"type:varchar(255);not null" json:"title"`
    Company               string    `gorm:"type:varchar(255)" json:"company"`
    Location              string    `gorm:"type:varchar(255)" json:"location"`
    PostedDate            string    `gorm:"type:varchar(50)" json:"posted_date"`
    Processed             bool      `json:"processed"`
    JobType               string    `gorm:"type:varchar(50)" json:"job_type"`
    Skills                string    `gorm:"type:text" json:"skills"`                        // Comma-separated list of required skills
    UserSkills            string    `gorm:"type:text" json:"user_skills"`                   // Comma-separated list of user skills
    MinSalary             int       `json:"min_salary"`                                    // Minimum salary
    MaxSalary             int       `json:"max_salary"`                                    // Maximum salary
    MatchScore            float64   `json:"match_score"`                                   // Match score from 0 to 100
    Description           string    `gorm:"type:longtext" json:"description"`               // Job description

    Selected              bool      `gorm:"default:false" json:"selected"`
    CvGenerated           bool      `gorm:"default:false" json:"cv_generated"`
    CoverLetterGenerated  bool      `gorm:"default:false" json:"cover_letter_generated"`
	ViewLink              bool  	`gorm:"default:false" json:"view_link"`

}


    type SeekerProfileDTO struct {
        ID                         uint      `json:"id"`
        AuthUserID                 uuid.UUID `json:"authUserId"`

        // From PersonalInfo
        FirstName   string  `json:"firstName"`
        SecondName  *string `json:"secondName,omitempty"`

        // From ProfessionalSummary
        Skills []string `json:"skills"`

        // From WorkExperience
        TotalExperienceInMonths int `json:"totalExperienceInMonths"`

        // From Certificate
        Certificates []string `json:"certificates"`

        // From PreferredJobTitle
        PreferredJobTitle string `json:"preferredJobTitle"`

        SubscriptionTier           string    `json:"subscriptionTier"`
        DailySelectableJobsCount   int       `json:"dailySelectableJobsCount"`
        DailyGeneratableCV         int       `json:"dailyGeneratableCv"`
        DailyGeneratableCoverletter int      `json:"dailyGeneratableCoverletter"`
        TotalApplications          int       `json:"totalApplications"`

        // New: Total number of jobs available across sources
        TotalJobsAvailable int `json:"totalJobsAvailable"`

        // New: Profile completion percentage
        ProfileCompletion int `json:"profileCompletion"`
    }

// LinkResponseDTO represents the response DTO for job application links
type LinkResponseDTO struct {
	JobID   string `json:"job_id"`
	JobLink string `json:"job_link"`
	Source  string `json:"source"`
}