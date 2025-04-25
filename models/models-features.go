package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// UserEntryTimeline for MongoDB
type UserEntryTimeline struct {
	ID                              primitive.ObjectID `bson:"_id,omitempty" json:"id"`  // MongoDB ID
	AuthUserID                       primitive.ObjectID `bson:"authUserId" json:"authUserId"`

	PersonalInfosCompleted          bool       `bson:"personalInfosCompleted" json:"personalInfosCompleted"`
	PersonalInfosRequired           bool       `bson:"personalInfosRequired" json:"personalInfosRequired"`

	ProfessionalSummariesCompleted  bool       `bson:"professionalSummariesCompleted" json:"professionalSummariesCompleted"`
	ProfessionalSummariesRequired   bool       `bson:"professionalSummariesRequired" json:"professionalSummariesRequired"`

	WorkExperiencesCompleted        bool       `bson:"workExperiencesCompleted" json:"workExperiencesCompleted"`
	WorkExperiencesRequired         bool       `bson:"workExperiencesRequired" json:"workExperiencesRequired"`

	EducationsCompleted             bool       `bson:"educationsCompleted" json:"educationsCompleted"`
	EducationsRequired              bool       `bson:"educationsRequired" json:"educationsRequired"`

	CertificatesCompleted           bool       `bson:"certificatesCompleted" json:"certificatesCompleted"`
	CertificatesRequired            bool       `bson:"certificatesRequired" json:"certificatesRequired"`

	LanguagesCompleted              bool       `bson:"languagesCompleted" json:"languagesCompleted"`
	LanguagesRequired               bool       `bson:"languagesRequired" json:"languagesRequired"`

	PreferredJobTitlesCompleted     bool       `bson:"preferredJobTitlesCompleted" json:"preferredJobTitlesCompleted"`
	PreferredJobTitlesRequired      bool       `bson:"preferredJobTitlesRequired" json:"preferredJobTitlesRequired"`

	Completed                       bool       `bson:"completed" json:"completed"`

	CreatedAt                       time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt                       time.Time  `bson:"updatedAt" json:"updatedAt"`
}

// SalaryRange for MongoDB
type SalaryRange struct {
	Min int `bson:"min" json:"min"`
	Max int `bson:"max" json:"max"`
}

// SelectedJobApplication for MongoDB
type SelectedJobApplication struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthUserID             primitive.ObjectID `bson:"authUserId" json:"authUserId"`
	Source                 string             `bson:"source" json:"source"`
	JobID                  string             `bson:"jobId" json:"jobId"`
	Title                 string             `bson:"title" json:"title"`
	Company               string             `bson:"company" json:"company"`
	Location              string             `bson:"location" json:"location"`
	PostedDate            string             `bson:"postedDate" json:"postedDate"`
	Processed             bool               `bson:"processed" json:"processed"`
	JobType               string             `bson:"jobType" json:"jobType"`
	Skills                string             `bson:"skills" json:"skills"`
	UserSkills            string             `bson:"userSkills" json:"userSkills"`
	MinSalary             int                `bson:"minSalary" json:"minSalary"`
	MaxSalary             int                `bson:"maxSalary" json:"maxSalary"`
	MatchScore            float64            `bson:"matchScore" json:"matchScore"`
	Description           string             `bson:"description" json:"description"`

	Selected              bool               `bson:"selected" json:"selected"`
	CvGenerated           bool               `bson:"cvGenerated" json:"cvGenerated"`
	CoverLetterGenerated  bool               `bson:"coverLetterGenerated" json:"coverLetterGenerated"`
	ViewLink              bool               `bson:"viewLink" json:"viewLink"`
}

// CV for MongoDB
type CV struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthUserID primitive.ObjectID `bson:"authUserId" json:"authUserId"`
	JobID      string             `bson:"jobId" json:"jobId"`
	CVUrl      string             `bson:"cvUrl" json:"cvUrl"`
}

// CoverLetter for MongoDB
type CoverLetter struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthUserID    primitive.ObjectID `bson:"authUserId" json:"authUserId"`
	JobID         string             `bson:"jobId" json:"jobId"`
	CoverLetterURL string            `bson:"coverLetterURL" json:"coverLetterURL"`
}
