package repositories

import (
    "RAAS/models"
    "gorm.io/gorm"
    "fmt"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// SeekerSignUpInput represents the incoming data for seeker signup
type SeekerSignUpInput struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Location  string `json:"location" binding:"required"`
}

// ValidateSeekerSignUpInput validates the Seeker signup input
func ValidateSeekerSignUpInput(input SeekerSignUpInput) error {
    if input.Email == "" || input.Password == "" || input.FirstName == "" || input.LastName == "" || input.Location == "" {
        return fmt.Errorf("all fields are required")
    }
    return nil
}

// IsEmailTaken checks if an email is already registered
func IsEmailTaken(db *gorm.DB, email string) (bool, error) {
    var count int64
    err := db.Model(&models.AuthUser{}).Where("email = ?", email).Count(&count).Error
    if err != nil {
        return false, err // Return false along with the actual error
    }
    return count > 0, nil // Returns true if email exists, false otherwise
}

// CreateSeeker creates a new seeker and their authentication record
func CreateSeeker(db *gorm.DB, input SeekerSignUpInput, hashedPassword string) error {
    // Create AuthUser
    authUser := models.AuthUser{
        Email:    input.Email,
        Password: hashedPassword,
        Role:     "seeker",
    }

    if err := db.Create(&authUser).Error; err != nil {
        return err
    }

    // Create Seeker
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

    return nil
}


// AuthenticateUser checks the credentials and returns the authenticated user.
func AuthenticateUser(db *gorm.DB, email, password string) (*models.AuthUser, error) {
    var user models.AuthUser

    if err := db.Where("email = ?", email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid email or password")
        }
        return nil, err
    }

    // Compare hashed password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid email or password")
    }

    return &user, nil
}
