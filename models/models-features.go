package models

import (
	"time"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

type UserEntryTimeline struct {
	ID                              primitive.ObjectID `bson:"_id,omitempty" json:"id"`  // MongoDB ID
	AuthUserID                       uuid.UUID         `bson:"auth_user_id" json:"auth_user_id"`

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

func CreateUserEntryTimelineIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "auth_user_id", Value: 1}},
		Options: options.Index().SetUnique(true),       
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}






type SalaryRange struct {
	Min int `bson:"min" json:"min"`
	Max int `bson:"max" json:"max"`
}
type SelectedJobApplication struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthUserID             uuid.UUID         `json:"auth_user_id" bson:"auth_user_id"`
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
func CreateSelectedJobApplicationIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "authUserId", Value: 1}, {Key: "jobId", Value: 1}}, // Compound index on authUserId and jobId
		Options: options.Index().SetUnique(true),                                 // Unique index
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}





type CV struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthUserID uuid.UUID         `json:"auth_user_id" bson:"auth_user_id"`
	JobID      string             `bson:"jobId" json:"jobId"`
	CVUrl      string             `bson:"cvUrl" json:"cvUrl"`
}

func CreateCVIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "authUserId", Value: 1}, {Key: "jobId", Value: 1}}, 
		Options: options.Index().SetUnique(true),                               
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}


type CoverLetter struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthUserID    uuid.UUID         `json:"auth_user_id" bson:"auth_user_id"`
	JobID         string             `bson:"jobId" json:"jobId"`
	CoverLetterURL string            `bson:"coverLetterURL" json:"coverLetterURL"`
}


func CreateCoverLetterIndexes(collection *mongo.Collection) error {

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "authUserId", Value: 1}, {Key: "jobId", Value: 1}}, 
		Options: options.Index().SetUnique(true),                                 
	}

	// Create the index
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
