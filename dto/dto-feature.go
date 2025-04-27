package dto

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Job Retrieval

type SalaryRange struct {
    Min int `json:"min" bson:"min"`
    Max int `json:"max" bson:"max"`
}

type JobDTO struct {
    Source         string       `json:"source" bson:"source"`                   // "linkedin" or "xing"
    ID             uint         `json:"id" bson:"id"`                           // UUID or unique DB ID
    JobID          string       `json:"job_id" bson:"job_id"`                   // Platform-specific Job ID
    Title          string       `json:"title" bson:"title"`
    Company        string       `json:"company" bson:"company"`
    Location       string       `json:"location" bson:"location"`
    PostedDate     string       `json:"posted_date" bson:"posted_date"`
    Processed      bool         `json:"processed" bson:"processed"`
    JobType        string       `json:"job_type" bson:"job_type"`               // e.g., Full-time, Part-time
    Skills         string       `json:"skills" bson:"skills"`                   // Comma-separated required skills
    UserSkills     []string     `json:"user_skills" bson:"user_skills"`         // List of user skills used in matching
    ExpectedSalary SalaryRange  `json:"expected_salary" bson:"expected_salary"` // Expected salary range
    MatchScore     float64      `json:"match_score" bson:"match_score"`         // Match score from 0 to 100
    Description    string       `json:"description" bson:"description"`         // Job description text
}

// JobFilterDTO represents the filter data for job retrieval.
type JobFilterDTO struct {
    Title string `form:"title" bson:"title"` // Query param: /jobs/linkedin?title=developer
}

// MatchScoreResponse represents the response containing the match score for a job.
type MatchScoreResponse struct {
    SeekerID   string  `json:"seeker_id" bson:"seeker_id"` // Changed to string
    JobID      string  `json:"job_id" bson:"job_id"`
    MatchScore float64 `json:"match_score" bson:"match_score"`
}

// SelectedJobApplication represents a job selected by a seeker.
type SelectedJobApplication struct {
    ID                   primitive.ObjectID `bson:"_id,omitempty"`
    AuthUserID           string             `bson:"auth_user_id"`            // Changed to string
    Source               string             `bson:"source"`                 // "linkedin" or "xing"
    JobID                string             `bson:"job_id"`                 // Platform-specific job ID
    Title                string             `bson:"title"`
    Company              string             `bson:"company"`
    Location             string             `bson:"location"`
    PostedDate           string             `bson:"posted_date"`
    Processed            bool               `bson:"processed"`
    JobType              string             `bson:"job_type"`
    Skills               string             `bson:"skills"`                 // Comma-separated list of required skills
    UserSkills           string             `bson:"user_skills"`            // Comma-separated list of user skills
    MinSalary            int                `bson:"min_salary"`             // Minimum salary
    MaxSalary            int                `bson:"max_salary"`             // Maximum salary
    MatchScore           float64            `bson:"match_score"`            // Match score from 0 to 100
    Description          string             `bson:"description"`           // Job description
    Selected             bool               `bson:"selected" default:"false"`
    CvGenerated          bool               `bson:"cv_generated" default:"false"`
    CoverLetterGenerated bool               `bson:"cover_letter_generated" default:"false"`
    ViewLink             bool               `bson:"view_link" default:"false"`
}

type SeekerProfileDTO struct {
    ID                         uint      `json:"id" bson:"id"`
    AuthUserID                 string    `json:"auth_user_id" bson:"auth_user_id"` // Changed to string

    // From PersonalInfo
    FirstName   string  `json:"first_name" bson:"first_name"`
    SecondName  *string `json:"second_name,omitempty" bson:"second_name,omitempty"`

    // From ProfessionalSummary
    Skills []string `json:"skills" bson:"skills"`

    // From WorkExperience
    TotalExperienceInMonths int `json:"total_experience_in_months" bson:"total_experience_in_months"`

    // From Certificate
    Certificates []string `json:"certificates" bson:"certificates"`

    // From PreferredJobTitle
    PreferredJobTitle string `json:"preferred_job_title" bson:"preferred_job_title"`

    SubscriptionTier           string    `json:"subscription_tier" bson:"subscription_tier"`
    DailySelectableJobsCount   int       `json:"daily_selectable_jobs_count" bson:"daily_selectable_jobs_count"`
    DailyGeneratableCV         int       `json:"daily_generatable_cv" bson:"daily_generatable_cv"`
    DailyGeneratableCoverletter int      `json:"daily_generatable_coverletter" bson:"daily_generatable_coverletter"`
    TotalApplications          int       `json:"total_applications" bson:"total_applications"`

    // New: Total number of jobs available across sources
    TotalJobsAvailable int `json:"total_jobs_available" bson:"total_jobs_available"`

    // New: Profile completion percentage
    ProfileCompletion int `json:"profile_completion" bson:"profile_completion"`
}

// LinkResponseDTO represents the response DTO for job application links
type LinkResponseDTO struct {
    JobID   string `json:"job_id" bson:"job_id"`
    JobLink string `json:"job_link" bson:"job_link"`
    Source  string `json:"source" bson:"source"`
}
