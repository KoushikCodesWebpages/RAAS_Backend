package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"fmt"
	"errors"

	"RAAS/config"
	"RAAS/models"
	"RAAS/dto"
	"RAAS/security"
	"RAAS/handlers/repo"
	"RAAS/utils"
)

func SeekerSignUp(c *gin.Context) {
	var input dto.SeekerSignUpInput

	// Bind input and check for errors
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Initialize userRepo with valid config
	db := c.MustGet("db").(*gorm.DB)
	userRepo := repo.NewUserRepo(db, config.Cfg)

	// Validate the input data
	if err := userRepo.ValidateSeekerSignUpInput(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email is already taken and whether it's verified or not
	var user models.AuthUser
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Database error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
			return
		}
	} else {
		// If email exists, check if it's verified
		if user.EmailVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already taken and verified."})
			return
		}


		// If email is not verified, resend the verification email
		token := user.VerificationToken
		verificationLink := fmt.Sprintf("%s/verify-email?token=%s", config.Cfg.FrontendBaseUrl, token)
		emailBody := fmt.Sprintf(`
			<p>Hello %s,</p>
			<p>Thanks for signing up! Please verify your email by clicking the link below:</p>
			<p><a href="%s">Verify Email</a></p>
			<p>If you did not sign up, you can ignore this email.</p>
		`, input.Email, verificationLink)

		// Log for debugging email sending
		fmt.Printf("Resending verification email to: %s\n", input.Email)

		// Send the verification email again
		emailCfg := utils.EmailConfig{
			Host:     config.Cfg.EmailHost,
			Port:     config.Cfg.EmailPort,
			Username: config.Cfg.EmailHostUser,
			Password: config.Cfg.EmailHostPassword,
			From:     config.Cfg.DefaultFromEmail,
			UseTLS:   config.Cfg.EmailUseTLS,
		}

		if err := utils.SendEmail(emailCfg, input.Email, "Verify your email", emailBody); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend verification email", "details": err.Error()})
			return
		}

		// Log after successful email sending
		fmt.Println("Verification email resent successfully.")

		c.JSON(http.StatusOK, gin.H{"message": "Please check your email to verify your account. Email resent."})
		return
	}

	// Log before hashing the password
	fmt.Println("Hashing password for user:", input.Email)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create the user and send the verification email
	fmt.Printf("Creating user with email: %s\n", input.Email)
	if err := userRepo.CreateSeeker(input, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Seeker", "details": err.Error()})
		return
	}

	// Log successful user creation
	fmt.Printf("User created successfully. Email verification sent to: %s\n", input.Email)

	// Final response
	c.JSON(http.StatusCreated, gin.H{"message": "Seeker registered successfully. Please check your email to verify."})
}


func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Missing token")
		return
	}

	db, exists := c.Get("db")
	if !exists {
		c.String(http.StatusInternalServerError, "DB not found")
		return
	}

	var user models.AuthUser
	if err := db.(*gorm.DB).Where("verification_token = ?", token).First(&user).Error; err != nil {
		c.String(http.StatusNotFound, "Invalid or expired token")
		return
	}

	// If the token is valid, verify the email
	if user.EmailVerified {
		c.String(http.StatusOK, "Email already verified.")
		return
	}

	// Mark the email as verified and clear the verification token
	user.EmailVerified = true
	user.VerificationToken = ""

	if err := db.(*gorm.DB).Save(&user).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to verify email")
		return
	}

	html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Email Verified</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f2f4f8;
					color: #333;
					text-align: center;
					padding-top: 100px;
				}
				.card {
					background: white;
					padding: 40px;
					margin: auto;
					border-radius: 8px;
					box-shadow: 0 4px 6px rgba(0,0,0,0.1);
					width: 90%;
					max-width: 500px;
				}
				h1 {
					color: #28a745;
				}
				p {
					margin-top: 10px;
					font-size: 18px;
				}
				a {
					display: inline-block;
					margin-top: 20px;
					text-decoration: none;
					color: white;
					background-color: #007bff;
					padding: 10px 20px;
					border-radius: 5px;
				}
			</style>
		</head>
		<body>
			<div class="card">
				<h1>✅ Email Verified</h1>
				<p>Your email has been successfully verified.</p>
				<a href="http://localhost:3000/login">Go to Login</a>
			</div>
		</body>
		</html>
	`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}




func Login(c *gin.Context) {
	var input dto.LoginInput
	db := c.MustGet("db").(*gorm.DB)
	userRepo := repo.NewUserRepo(db, config.Cfg)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	user, err := userRepo.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := security.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}


