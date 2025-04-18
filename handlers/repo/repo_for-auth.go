package repo

import (
	"RAAS/dto"
	"RAAS/models"
	"fmt"
	"gorm.io/gorm"
	"RAAS/utils"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"RAAS/config"
)

type UserRepo struct {
	DB     *gorm.DB

}

func NewUserRepo(db *gorm.DB, config *config.Config) *UserRepo {
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
		Phone:			   input.Number,
		Role:              "seeker",
		VerificationToken: token,
		EmailVerified:     false, // Assume false until verified
	}

	// Save AuthUser to the database
	if err := r.DB.Create(&authUser).Error; err != nil {
		return fmt.Errorf("failed to create auth user: %v", err)
	}

	// Create associated Seeker profile
	seeker := models.Seeker{
		AuthUserID:       authUser.ID,
		SubscriptionTier: "free",
	}

	// Save Seeker to the database
	if err := r.DB.Create(&seeker).Error; err != nil {
		return fmt.Errorf("failed to create seeker profile: %v", err)
	}

	// Create UserEntryTimeline for the new user
	timeline := models.UserEntryTimeline{
		UserID:                         authUser.ID,
	}
	if err := r.DB.Create(&timeline).Error; err != nil {
		return fmt.Errorf("user created but failed to create entry timeline: %v", err)
	}

	// Construct email verification link (optional)
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", config.Cfg.FrontendBaseUrl, token)

	emailBody := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Thanks for signing up! Please verify your email by clicking the link below:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>If you did not sign up, you can ignore this email.</p>
	`, input.Email, verificationLink)

	// Prepare email config from loaded app config
	emailCfg := utils.EmailConfig{
		Host:     config.Cfg.EmailHost,
		Port:     config.Cfg.EmailPort,
		Username: config.Cfg.EmailHostUser,
		Password: config.Cfg.EmailHostPassword,
		From:     config.Cfg.DefaultFromEmail,
		UseTLS:   config.Cfg.EmailUseTLS,
	}

	// Send the verification email
	if err := utils.SendEmail(emailCfg, input.Email, "Verify your email", emailBody); err != nil {
		return fmt.Errorf("user created but failed to send verification email: %v", err)
	}

	return nil
}
func (r *UserRepo) AuthenticateUser(email, password string) (*models.AuthUser, error) {
    var user models.AuthUser

    // Find the user by email
    if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid email or password")
        }
        return nil, err
    }

    // Check if the email is verified
    if !user.EmailVerified {
        return nil, errors.New("email is not verified")
    }

    // Compare the provided password with the stored hashed password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid email or password")
    }

    return &user, nil
}
