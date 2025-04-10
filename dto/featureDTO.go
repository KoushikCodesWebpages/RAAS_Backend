package dto

import (
    // "RAAS/models"
	"github.com/google/uuid"
	"time"
)

//JOB TITLE

type JobTitleInput struct {
	PrimaryTitle   string  `json:"primaryTitle"`
	SecondaryTitle *string `json:"secondaryTitle,omitempty"`
	TertiaryTitle  *string `json:"tertiaryTitle,omitempty"`
}

//JOB RETRIEVAL

type JobDTO struct {
	Source     string `json:"source"`      // "linkedin" or "xing"
	ID         string `json:"id"`
	JobID      string `json:"job_id"`
	Title      string `json:"title"`
	Company    string `json:"company"`
	Location   string `json:"location"`
	PostedDate string `json:"posted_date"`
	Processed  bool   `json:"processed"`
}


//JOB FILTER

type JobFilterDTO struct {
	Title string `form:"title"` // passed as query param, e.g. /jobs/linkedin?title=developer
}

//JOB MATCH SCORE

type MatchScoreRequest struct {
	JobID  string `json:"job_id" binding:"required"`
}

type MatchScoreResponse struct {
	Score     float64   `json:"score"`
	Source    string    `json:"source"` // "cached" or "calculated"
	JobID     string    `json:"job_id"`
	UserID    uuid.UUID `json:"user_id"`
	Platform  string    `json:"platform,omitempty"`
	MatchedAt time.Time `json:"matched_at,omitempty"`
}
