package models

import (
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthUser struct {
	AuthUserID           string     `json:"auth_user_id" bson:"auth_user_id"` // Changed to string for MongoDB UUID storage
	Email                string     `json:"email" bson:"email"`
	Phone                string     `json:"phone" bson:"phone"`
	Password             string     `json:"password" bson:"password"`
	Role                 string     `json:"role" bson:"role"`
	EmailVerified        bool       `json:"email_verified" bson:"email_verified"`
	Provider             string     `json:"provider" bson:"provider,omitempty"`
	ResetTokenExpiry     *time.Time `json:"reset_token_expiry" bson:"reset_token_expiry"`
	IsActive             bool       `json:"is_active" bson:"is_active"`
	VerificationToken    string     `json:"verification_token" bson:"verification_token"`
	CreatedBy            string     `json:"created_by" bson:"created_by"` // Changed to string for MongoDB UUID storage
	UpdatedBy            string     `json:"updated_by" bson:"updated_by"` // Changed to string for MongoDB UUID storage
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


type Seeker struct {
	ID                          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AuthUserID                  string             `json:"auth_user_id" bson:"auth_user_id"`
	SubscriptionTier            string             `json:"subscription_tier" bson:"subscription_tier"`
	DailySelectableJobsCount    int                `json:"daily_selectable_jobs_count" bson:"daily_selectable_jobs_count"`
	DailyGeneratableCV          int                `json:"daily_generatable_cv" bson:"daily_generatable_cv"`
	DailyGeneratableCoverletter int                `json:"daily_generatable_coverletter" bson:"daily_generatable_coverletter"`
	TotalApplications           int                `json:"total_applications" bson:"total_applications"`

	PersonalInfo                bson.M             `json:"personal_info" bson:"personal_info"`
	ProfessionalSummary         bson.M             `json:"professional_summary" bson:"professional_summary"`
	
	WorkExperiences             []bson.M            `json:"work_experiences" bson:"work_experiences"`
	Education                 	[]bson.M            `json:"education" bson:"education"`
	Certificates                []bson.M            `json:"certificates" bson:"certificates"`
	Languages                   []bson.M            `json:"languages" bson:"languages"`

	PrimaryTitle                string             `json:"primary_title" bson:"primary_title"`
	SecondaryTitle              *string            `json:"secondary_title,omitempty" bson:"secondary_title,omitempty"`
	TertiaryTitle               *string            `json:"tertiary_title,omitempty" bson:"tertiary_title,omitempty"`
}

func CreateSeekerIndexes(collection *mongo.Collection) error {
	// Create index for AuthUserID to be unique
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "auth_user_id", Value: 1}}, 
		Options: options.Index().SetUnique(true),        
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

type Admin struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AuthUserID string             `json:"auth_user_id" bson:"auth_user_id"` // Change uuid.UUID to string
}

func CreateAdminIndexes(collection *mongo.Collection) error {
	// Create index for AuthUserID to be unique
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "auth_user_id", Value: 1}}, 
		Options: options.Index().SetUnique(true),      
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
