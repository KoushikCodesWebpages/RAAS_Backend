package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

)

// AUTH MODELS

type AuthUser struct {
    gorm.Model
    ID                   uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
    Email                string     `gorm:"unique;not null" json:"email"`
    Phone                string     `gorm:"not null" json:"phone"`
    Password             string     `json:"password"`
    Role                 string     `json:"role"`
    VerificationToken    string     `json:"verification_token"`
    EmailVerified        bool       `json:"email_verified"`
    Provider             string     `gorm:"default:'local'" json:"provider"`
    ResetTokenExpiry     *time.Time `json:"reset_token_expiry"` 

    IsActive             bool       `gorm:"default:true" json:"is_active"`

    CreatedBy            uuid.UUID `json:"created_by"`
    UpdatedBy            uuid.UUID `json:"updated_by"`
    DeletedAt            *time.Time `json:"deleted_at,omitempty"`

	
    LastLoginAt          *time.Time `json:"last_login_at,omitempty"`
    PasswordLastUpdated  *time.Time `json:"password_last_updated,omitempty"`
    TwoFactorEnabled     bool       `gorm:"default:false" json:"two_factor_enabled"`
    TwoFactorSecret      *string    `json:"two_factor_secret,omitempty"`

}


// SEEKER
type Seeker struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);uniqueIndex;not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"authUserId"`


    //SERVICE
	SubscriptionTier           string         `gorm:"default:'free'" json:"subscriptionTier"`
	DailySelectableJobsCount   int            `gorm:"default:5" json:"dailySelectableJobsCount"`
	DailyGeneratableCV         int            `gorm:"default:100" json:"dailyGeneratableCv"`
	DailyGeneratableCoverletter int           `gorm:"default:100" json:"dailyGeneratableCoverletter"`
	TotalApplications          int            `gorm:"default:0" json:"totalApplications"`

	//DATA
	PersonalInfo         datatypes.JSON `gorm:"type:json" json:"personalInfo"`
	ProfessionalSummary  datatypes.JSON `gorm:"type:json" json:"professionalSummary"`
	WorkExperiences      datatypes.JSON `gorm:"type:json" json:"workExperiences"`
	Educations           datatypes.JSON `gorm:"type:json" json:"education"`
	Certificates         datatypes.JSON `gorm:"type:json" json:"certificates"`
	Languages            datatypes.JSON `gorm:"type:json" json:"languages"`

	// JOB TITLES
	PrimaryTitle   string  `gorm:"type:varchar(255);" json:"primaryTitle"`
	SecondaryTitle *string `gorm:"type:varchar(255);" json:"secondaryTitle"`
	TertiaryTitle  *string `gorm:"type:varchar(255);" json:"tertiaryTitle"`
}

// ADMIN

type Admin struct {
	gorm.Model
	AuthUserID uuid.UUID `gorm:"type:char(36);unique;not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"authUserId"`
}
