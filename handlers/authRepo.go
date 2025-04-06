package handlers

import (
	"RAAS/dto"
	"RAAS/models"
	//"RAAS/security"
	"fmt"

	"gorm.io/gorm"

	//"RAAS/utils"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"RAAS/config"
)
func validateSeekerSignUpInput(input dto.SeekerSignUpInput) error {
    if input.Email == "" || input.Password == "" || input.FirstName == "" || input.LastName == "" || input.Location == "" {
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
    // Generate a verification token
    token := uuid.New().String()

    // Create AuthUser with verification fields
    authUser := models.AuthUser{
        Email:             input.Email,
        Password:          hashedPassword,
        Role:              "seeker",
        VerificationToken: token,
        //for now true
        EmailVerified:     true,
    }

    // Save AuthUser to the database
    if err := db.Create(&authUser).Error; err != nil {
        return err
    }

    // Create associated Seeker profile
    seeker := models.Seeker{
        AuthUserID: authUser.ID,
        AuthUser:   authUser,
        FirstName:  input.FirstName,
        LastName:   input.LastName,
        Location:   input.Location,
    }

    if err := db.Create(&seeker).Error; err != nil {
        return err
    }

    // Construct email verification link
    /*verificationLink := fmt.Sprintf("https://your-frontend.com/verify-email?token=%s", token)
    emailBody := fmt.Sprintf(`
        <p>Hello %s,</p>
        <p>Thanks for signing up! Please verify your email by clicking the link below:</p>
        <p><a href="%s">Verify Email</a></p>
        <p>If you did not sign up, you can ignore this email.</p>
    `, input.FirstName, verificationLink)

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
        */

    return nil
}

// func createSeekerFromGoogleOAuth(db *gorm.DB, info security.UserInfo) (*models.AuthUser, error) {
//     var user models.AuthUser

//     // Check if user exists already via Google
//     result := db.Where("email = ? AND provider = ?", info.Email, "google").First(&user)
//     if result.Error == nil {
//         return &user, nil
//     }

//     // Create AuthUser
//     authUser := models.AuthUser{
//         Email:         info.Email,
//         Role:          "seeker",
//         Provider:      "google",
//         EmailVerified: true,
//     }

//     if err := db.Create(&authUser).Error; err != nil {
//         return nil, err
//     }

//     // Create Seeker profile (you can refine name splitting later)
//     seeker := models.Seeker{
//         AuthUserID: authUser.ID,
//         AuthUser:   authUser,
//         FirstName:  info.Name,
//         LastName:   "",
//         Location:   "",
//     }

//     if err := db.Create(&seeker).Error; err != nil {
//         return nil, err
//     }

//     return &authUser, nil
// }


func authenticateUser(db *gorm.DB, email, password string) (*models.AuthUser, error) {
    var user models.AuthUser

    if err := db.Where("email = ?", email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid email or password")
        }
        return nil, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid email or password")
    }

    return &user, nil
}
