package dto

import (
    // "RAAS/models"
	"github.com/google/uuid"
	"time"
)


//JOB RETRIEVAL

// JobDTO represents the job data to be returned in the job listing.
type JobDTO struct {
	Source     string  `json:"source"`      // "linkedin" or "xing"
	ID         string  `json:"id"`
	JobID      string  `json:"job_id"`
	Title      string  `json:"title"`
	Company    string  `json:"company"`
	Location   string  `json:"location"`
	PostedDate string  `json:"posted_date"`
	Processed  bool    `json:"processed"`
	JobType    string  `json:"job_type,omitempty"`  // Job type (e.g., Full-time, Part-time)
	Skills     string  `json:"skills,omitempty"`    // Comma-separated skills required for the job
	MatchScore float64 `json:"match_score"`         // Match score between user and job (from 0 to 100)
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
