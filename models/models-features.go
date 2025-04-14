package models

import (
    "github.com/google/uuid"
	//"gorm.io/datatypes"
    "gorm.io/gorm"
	"time"
)

type UserEntryTimeline struct {
    ID                              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
    UserID                          uuid.UUID  `gorm:"type:char(36);uniqueIndex;not null" json:"user_id"`

    PersonalInfosCompleted          bool       `gorm:"default:false" json:"personal_infos_completed"`
    PersonalInfosRequired           bool       `gorm:"default:true" json:"personal_infos_required"`

    ProfessionalSummariesCompleted  bool       `gorm:"default:false" json:"professional_summaries_completed"`
    ProfessionalSummariesRequired   bool       `gorm:"default:true" json:"professional_summaries_required"`

    WorkExperiencesCompleted        bool       `gorm:"default:false" json:"work_experiences_completed"`
    WorkExperiencesRequired         bool       `gorm:"default:true" json:"work_experiences_required"`

    EducationsCompleted             bool       `gorm:"default:false" json:"educations_completed"`
    EducationsRequired              bool       `gorm:"default:true" json:"educations_required"`

    CertificatesCompleted           bool       `gorm:"default:false" json:"certificates_completed"`
    CertificatesRequired            bool       `gorm:"default:false" json:"certificates_required"`

    LanguagesCompleted              bool       `gorm:"default:false" json:"languages_completed"`
    LanguagesRequired               bool       `gorm:"default:false" json:"languages_required"`

    PreferredJobTitlesCompleted     bool       `gorm:"default:false" json:"preferred_job_titles_completed"`
    PreferredJobTitlesRequired      bool       `gorm:"default:true" json:"preferred_job_titles_required"`

    Completed                       bool       `gorm:"default:false" json:"completed"` // ✅ New Field

    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}


type SalaryRange struct {
    Min int `json:"min"`
    Max int `json:"max"`
}


type SelectedJobApplication struct {
	gorm.Model
	AuthUserID            uuid.UUID `gorm:"type:char(36);not null" json:"authUserId"`       // Foreign key to Seeker
	Source                string    `gorm:"type:varchar(20);not null" json:"source"`        // "linkedin" or "xing"
	JobID                 string    `gorm:"type:varchar(100);not null" json:"job_id"`       // Platform-specific job ID
	Title                 string    `gorm:"type:varchar(255);not null" json:"title"`
	Company               string    `gorm:"type:varchar(255)" json:"company"`
	Location              string    `gorm:"type:varchar(255)" json:"location"`
	PostedDate            string    `gorm:"type:varchar(50)" json:"posted_date"`
	Processed             bool      `json:"processed"`
	JobType               string    `gorm:"type:varchar(50)" json:"job_type"`
	Skills                string    `gorm:"type:text" json:"skills"`                        // Comma-separated
	UserSkills            string    `gorm:"type:text" json:"user_skills"`                   // Comma-separated
	MinSalary             int       `json:"min_salary"`
	MaxSalary             int       `json:"max_salary"`
	MatchScore            float64   `json:"match_score"`
	Description           string    `gorm:"type:longtext" json:"description"`

	Selected              bool      `gorm:"default:false" json:"selected"`
	CvGenerated           bool      `gorm:"default:false" json:"cv_generated"`
	CoverLetterGenerated  bool      `gorm:"default:false" json:"cover_letter_generated"`
}

type MatchScore struct {
	SeekerID   uuid.UUID `gorm:"type:char(36);primaryKey"`  // Seeker ID (foreign key reference)
	JobID      string    `gorm:"primaryKey"`                // Job ID (foreign key reference, as string)
	MatchScore float64   `gorm:"type:float"`                // Match score percentage (0 to 100)
}