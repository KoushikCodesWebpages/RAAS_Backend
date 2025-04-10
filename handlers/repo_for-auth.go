package handlers

import (
	"RAAS/dto"
	"RAAS/models"
	//"RAAS/security"
	"fmt"

	"gorm.io/gorm"

	"RAAS/utils"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"RAAS/config"
)

var appConfig *config.Config

func init() {
	var err error
	appConfig, err = config.InitConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
}


func validateSeekerSignUpInput(input dto.SeekerSignUpInput) error {
    if input.Email == "" || input.Password == "" || input.Number == "" {
        return fmt.Errorf("all fields are required")
    }
    return nil
}

func isEmailTaken(db *gorm.DB, email string) (bool, error) {
    var count int64
    err := db.Model(&models.AuthUser{}).Where("email = ?", email).Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func createSeeker(db *gorm.DB, input dto.SeekerSignUpInput, hashedPassword string, config *config.Config) error {
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
	if err := db.Create(&authUser).Error; err != nil {
		return fmt.Errorf("failed to create auth user: %v", err)
	}

	// Create associated Seeker profile, using the correct AuthUserID
	seeker := models.Seeker{
		AuthUserID: authUser.ID,
		SubscriptionTier: "free",
	    // Link Seeker to AuthUser by ID
		   // Assign Location
	}

	// Save Seeker to the database
	if err := db.Create(&seeker).Error; err != nil {
		return fmt.Errorf("failed to create seeker profile: %v", err)
	}

	// Construct email verification link (optional)
	fmt.Println("Loaded frontend base URL:", appConfig.FrontendBaseUrl)
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", config.FrontendBaseUrl, token)

	emailBody := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Thanks for signing up! Please verify your email by clicking the link below:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>If you did not sign up, you can ignore this email.</p>
	`, input.Email, verificationLink)
	

	// Prepare email config from loaded app config
	emailCfg := utils.EmailConfig{
		Host:     config.EmailHost,
		Port:     config.EmailPort,
		Username: config.EmailHostUser,
		Password: config.EmailHostPassword,
		From:     config.DefaultFromEmail,
		UseTLS:   config.EmailUseTLS,
	}

	// Send the verification email
	if err := utils.SendEmail(emailCfg, input.Email, "Verify your email", emailBody); err != nil {
		return fmt.Errorf("user created but failed to send verification email: %v", err)
	}
	

	return nil
}


func authenticateUser(db *gorm.DB, email, password string) (*models.AuthUser, error) {
    var user models.AuthUser

    // Find the user by email
    if err := db.Where("email = ?", email).First(&user).Error; err != nil {
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

