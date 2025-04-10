package dto

import (
    // "RAAS/models"
	// "github.com/google/uuid"
)


type JobTitleInput struct {
	PrimaryTitle   string  `json:"primaryTitle"`
	SecondaryTitle *string `json:"secondaryTitle,omitempty"`
	TertiaryTitle  *string `json:"tertiaryTitle,omitempty"`
}

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

type JobFilterDTO struct {
	Title string `form:"title"` // passed as query param, e.g. /jobs/linkedin?title=developer
}
