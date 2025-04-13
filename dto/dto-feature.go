package dto

import (
    // "RAAS/models"
	"github.com/google/uuid"
	"time"
)


//JOB RETRIEVAL


type SalaryRange struct {
    Min int `json:"min"`
    Max int `json:"max"`
}
type JobDTO struct {
    Source         string       `json:"source"`                   // "linkedin" or "xing"
    ID             string       `json:"id"`                       // UUID or unique DB ID
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

// MatchScoreRequest represents the request to calculate a match score.
type MatchScoreRequest struct {
	JobID string `json:"job_id" binding:"required"`
}

// MatchScoreResponse represents the response containing the match score for a job.
type MatchScoreResponse struct {
	Score     float64   `json:"score"`
	Source    string    `json:"source"`    // "cached" or "calculated"
	JobID     string    `json:"job_id"`
	UserID    uuid.UUID `json:"user_id"`
	Platform  string    `json:"platform,omitempty"` // Platform (linkedin/xing)
	MatchedAt time.Time `json:"matched_at,omitempty"`
}
