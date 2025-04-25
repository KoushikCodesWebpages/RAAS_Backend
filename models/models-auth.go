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

// AUTH MODELS
type AuthUser struct {
	AuthUserID           uuid.UUID  `json:"auth_user_id" bson:"auth_user_id,omitempty"`
	Email                string     `json:"email" bson:"email"`
	Phone                string     `json:"phone" bson:"phone"`
	Password             string     `json:"password" bson:"password"`
	Role                 string     `json:"role" bson:"role"`
	EmailVerified        bool       `json:"email_verified" bson:"email_verified"`
	Provider             string     `json:"provider" bson:"provider,omitempty"`
	ResetTokenExpiry     *time.Time `json:"reset_token_expiry" bson:"reset_token_expiry"`
	IsActive             bool       `json:"is_active" bson:"is_active"`
	VerificationToken    string     `json:"verification_token" bson:"verification_token"`
	CreatedBy            uuid.UUID  `json:"created_by" bson:"created_by"`
	UpdatedBy            uuid.UUID  `json:"updated_by" bson:"updated_by"`
	LastLoginAt          *time.Time `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
	PasswordLastUpdated  *time.Time `json:"password_last_updated,omitempty" bson:"password_last_updated,omitempty"`
	TwoFactorEnabled     bool       `json:"two_factor_enabled" bson:"two_factor_enabled"`
	TwoFactorSecret      *string    `json:"two_factor_secret,omitempty" bson:"two_factor_secret,omitempty"`
}


func CreateAuthUserIndexes(collection *mongo.Collection) error {
	indexModelEmail := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	indexModelPhone := mongo.IndexModel{
		Keys:    bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	indexModelCompound := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}, {Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		indexModelEmail,
		indexModelPhone,
		indexModelCompound,
	})
	return err
}



// SEEKER
type Seeker struct {
	ID                      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AuthUserID              uuid.UUID         `json:"authUserId" bson:"auth_user_id"`
	SubscriptionTier        string            `json:"subscriptionTier" bson:"subscription_tier"`
	DailySelectableJobsCount int               `json:"dailySelectableJobsCount" bson:"daily_selectable_jobs_count"`
	DailyGeneratableCV      int               `json:"dailyGeneratableCv" bson:"daily_generatable_cv"`
	DailyGeneratableCoverletter int           `json:"dailyGeneratableCoverletter" bson:"daily_generatable_coverletter"`
	TotalApplications       int               `json:"totalApplications" bson:"total_applications"`

	// Embedded documents
	PersonalInfo            interface{}       `json:"personalInfo" bson:"personal_info"`
	ProfessionalSummary     interface{}       `json:"professionalSummary" bson:"professional_summary"`
	WorkExperiences         interface{}       `json:"workExperiences" bson:"work_experiences"`
	Educations              interface{}       `json:"education" bson:"education"`
	Certificates            interface{}       `json:"certificates" bson:"certificates"`
	Languages               interface{}       `json:"languages" bson:"languages"`

	// Job Titles
	PrimaryTitle           string            `json:"primaryTitle" bson:"primary_title"`
	SecondaryTitle         *string           `json:"secondaryTitle,omitempty" bson:"secondary_title,omitempty"`
	TertiaryTitle          *string           `json:"tertiaryTitle,omitempty" bson:"tertiary_title,omitempty"`
}

func CreateSeekerIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "auth_user_id", Value: 1}}, 
		Options: options.Index().SetUnique(true),        
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}


type Admin struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AuthUserID uuid.UUID         `json:"authUserId" bson:"auth_user_id"`
}

func CreateAdminIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "auth_user_id", Value: 1}}, 
		Options: options.Index().SetUnique(true),      
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}