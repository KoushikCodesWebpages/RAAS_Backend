package models

import (

	"gorm.io/datatypes"
)
type SalaryRange struct {
    Min int `json:"min"`
    Max int `json:"max"`
}

type Job struct {
    Source        string       `gorm:"not null" json:"source"` // "linkedin" or "xing"
    ID            string       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
    JobID         string       `gorm:"not null" json:"job_id"`
    Title         string       `gorm:"not null" json:"title"`
    Company       string       `gorm:"not null" json:"company"`
    Location      string       `gorm:"not null" json:"location"`
    PostedDate    string       `gorm:"not null" json:"posted_date"`
    Processed     bool         `gorm:"default:false" json:"processed"`
    JobType       string       `gorm:"nullable" json:"job_type,omitempty"`
    Skills        datatypes.JSON `gorm:"type:json;not null" json:"skills"` // Stored as JSON
    UserSkills    string       `gorm:"nullable" json:"user_skills,omitempty"`
    ExpectedSalary SalaryRange `gorm:"nullable" json:"expected_salary,omitempty"`
    MatchScore    float64      `gorm:"not null" json:"match_score"`
}