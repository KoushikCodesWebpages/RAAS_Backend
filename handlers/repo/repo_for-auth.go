package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"RAAS/config"
	"RAAS/dto"
	"RAAS/models"
	"RAAS/utils"
)

type UserRepo struct {
	DB *mongo.Database
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) ValidateSeekerSignUpInput(input dto.SeekerSignUpInput) error {
	if input.Email == "" || input.Password == "" || input.Number == "" {
		return fmt.Errorf("all fields are required")
	}
	return nil
}

func (r *UserRepo) CheckDuplicateEmailOrPhone(email, phone string) (bool, bool, error) {
	var user models.AuthUser

	filter := bson.M{
		"$or": []bson.M{
			{"email": email},
			{"phone": phone},
		},
	}

	err := r.DB.Collection("auth_users").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, false, nil
		}
		return false, false, fmt.Errorf("failed to check email or phone: %w", err)
	}

	emailExists := user.Email == email
	phoneExists := user.Phone == phone

	return emailExists, phoneExists, nil
}

func (r *UserRepo) CreateSeeker(input dto.SeekerSignUpInput, hashedPassword string) error {
	// Check for duplicate email or phone
	emailTaken, phoneTaken, err := r.CheckDuplicateEmailOrPhone(input.Email, input.Number)
	if err != nil {
		return fmt.Errorf("error checking for duplicates: %w", err)
	}
	if emailTaken {
		return fmt.Errorf("email is already taken")
	}
	if phoneTaken {
		return fmt.Errorf("phone number is already taken")
	}

	authUserID := uuid.New().String()
	token := uuid.New().String()

	// Create AuthUser
	authUser := models.AuthUser{
		AuthUserID:        authUserID,
		Email:             input.Email,
		Password:          hashedPassword,
		Phone:             input.Number,
		Role:              "seeker",
		EmailVerified:     false,
		VerificationToken: token,
		IsActive:          true,
		CreatedBy:         authUserID,
		UpdatedBy:         authUserID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert AuthUser
	_, err = r.DB.Collection("auth_users").InsertOne(ctx, authUser)
	if err != nil {
		return fmt.Errorf("failed to create auth user: %w", err)
	}

	// Create Seeker
	seeker := models.Seeker{
		AuthUserID:                  authUserID,
		SubscriptionTier:            "free",
		DailySelectableJobsCount:    10,
		DailyGeneratableCV:          100,
		DailyGeneratableCoverletter: 100,
		TotalApplications:           0,
		PersonalInfo:                bson.M{},
		ProfessionalSummary:         bson.M{},
		WorkExperiences:             []bson.M{},
		Education:                  []bson.M{},
		Certificates:                []bson.M{},
		Languages:                   []bson.M{},
		PrimaryTitle:                "",
		SecondaryTitle:              nil,
		TertiaryTitle:               nil,
	}

	_, err = r.DB.Collection("seekers").InsertOne(ctx, seeker)
	if err != nil {
		return fmt.Errorf("failed to create seeker profile: %w", err)
	}

	// Create Timeline
	timeline := models.UserEntryTimeline{
		AuthUserID:                     authUserID,
		PersonalInfosCompleted:         false,
		PersonalInfosRequired:          true,
		ProfessionalSummariesCompleted: false,
		ProfessionalSummariesRequired:  true,

		WorkExperiencesCompleted:       false,
		WorkExperiencesRequired:        false,
		EducationsCompleted:            false,
		EducationsRequired:             false,
		CertificatesCompleted:          false,
		CertificatesRequired:           false,
		LanguagesCompleted:             false,
		LanguagesRequired:              false,
		
		PreferredJobTitlesCompleted:    false,
		PreferredJobTitlesRequired:     true,
		Completed:                      false,
		CreatedAt:                      time.Now(),
		UpdatedAt:                      time.Now(),
	}

	_, err = r.DB.Collection("user_entry_timelines").InsertOne(ctx, timeline)
	if err != nil {
		return fmt.Errorf("user created but failed to create entry timeline: %w", err)
	}

	// Prepare verification email
	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", config.Cfg.Project.FrontendBaseUrl, token)
	emailBody := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Thanks for signing up! Please verify your email by clicking the link below:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>If you did not sign up, you can ignore this email.</p>
	`, input.Email, verificationLink)

	emailCfg := utils.EmailConfig{
		Host:     config.Cfg.Cloud.EmailHost,
		Port:     config.Cfg.Cloud.EmailPort,
		Username: config.Cfg.Cloud.EmailHostUser,
		Password: config.Cfg.Cloud.EmailHostPassword,
		From:     config.Cfg.Cloud.DefaultFromEmail,
		UseTLS:   config.Cfg.Cloud.EmailUseTLS,
	}

	if err := utils.SendEmail(emailCfg, input.Email, "Verify your email", emailBody); err != nil {
		return fmt.Errorf("user created but failed to send verification email: %w", err)
	}

	return nil
}

func (r *UserRepo) AuthenticateUser(ctx context.Context, email, password string) (*models.AuthUser, error) {
	var user models.AuthUser

	err := r.DB.Collection("auth_users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("incorrect password")
	}

	return &user, nil
}
