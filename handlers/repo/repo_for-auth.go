package repo

import (
	"fmt"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"RAAS/config"
	"RAAS/utils"
	"RAAS/dto"
	"RAAS/models"
)

type UserRepo struct {
	DB     *mongo.Client
}

func NewUserRepo(db *mongo.Client, config *config.Config) *UserRepo {
	return &UserRepo{
		DB:     db,
	}
}

func (r *UserRepo) ValidateSeekerSignUpInput(input dto.SeekerSignUpInput) error {
    if input.Email == "" || input.Password == "" || input.Number == "" {
        return fmt.Errorf("all fields are required")
    }
    return nil
}

func (r *UserRepo) CreateSeeker(input dto.SeekerSignUpInput, hashedPassword string) error {
	// Generate a verification token (UUID)
	token := uuid.New().String()

	// Create AuthUser with verification fields
	authUser := models.AuthUser{
		ID:                uuid.New(), // Generate UUID for AuthUser
		Email:             input.Email,
		Password:          hashedPassword,
		Phone:             input.Number,
		Role:              "seeker",
		VerificationToken: token,
		EmailVerified:     false, // Assume false until verified
	}

	// Save AuthUser to the database
	_, err := r.DB.Database(config.Cfg.Cloud.MongoDBName).Collection("auth_users").InsertOne(nil, authUser)
	if err != nil {
		return fmt.Errorf("failed to create auth user: %w", err)
	}

	// Create associated Seeker profile with default values
	seeker := models.Seeker{
		AuthUserID:                authUser.ID,
		SubscriptionTier:          "free", // Default value for subscription tier
		DailySelectableJobsCount:  5,     // Default value
		DailyGeneratableCV:       100,    // Default value
		DailyGeneratableCoverletter: 100, // Default value
		TotalApplications:        0,      // Default value
		PersonalInfo:             nil,    // or initialize with an empty JSON object
		ProfessionalSummary:      nil,    // or initialize with an empty JSON object
		WorkExperiences:          nil,    // or initialize with an empty JSON object
		PrimaryTitle:             "",     // You can leave it empty initially
		SecondaryTitle:           nil,    // You can leave it nil initially
		TertiaryTitle:            nil,    // You can leave it nil initially
	}

	// Save Seeker to the database
	_, err = r.DB.Database(config.Cfg.Cloud.MongoDBName).Collection("seekers").InsertOne(nil, seeker)
	if err != nil {
		return fmt.Errorf("failed to create seeker profile: %w", err)
	}

	timeline := models.UserEntryTimeline{
		AuthUserID: authUser.ID.String(), // Convert UUID to string
	}
	
	_, err = r.DB.Database(config.Cfg.Cloud.MongoDBName).Collection("user_entry_timelines").InsertOne(nil, timeline)
	if err != nil {
		return fmt.Errorf("user created but failed to create entry timeline: %w", err)
	}

	// Construct email verification link
	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", config.Cfg.Project.FrontendBaseUrl, token)

	// Create email body
	emailBody := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Thanks for signing up! Please verify your email by clicking the link below:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>If you did not sign up, you can ignore this email.</p>
	`, input.Email, verificationLink)

	// Prepare email config from loaded app config
	emailCfg := utils.EmailConfig{
		Host:     config.Cfg.Cloud.EmailHost,
		Port:     config.Cfg.Cloud.EmailPort,
		Username: config.Cfg.Cloud.EmailHostUser,
		Password: config.Cfg.Cloud.EmailHostPassword,
		From:     config.Cfg.Cloud.DefaultFromEmail,
		UseTLS:   config.Cfg.Cloud.EmailUseTLS,
	}

	// Send the verification email
	if err := utils.SendEmail(emailCfg, input.Email, "Verify your email", emailBody); err != nil {
		return fmt.Errorf("user created but failed to send verification email: %w", err)
	}

	return nil
}

func (r *UserRepo) AuthenticateUser(email, password string) (*models.AuthUser, error) {
	var user models.AuthUser

	// Find the user by email
	err := r.DB.Database(config.Cfg.Cloud.MongoDBName).Collection("auth_users").FindOne(nil, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}

	// Check if the password matches
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("incorrect password")
	}

	// Return the authenticated user
	return &user, nil
}
