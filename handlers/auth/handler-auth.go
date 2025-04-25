package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"


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
	db := c.MustGet("db").(*mongo.Client)
	userRepo := repo.NewUserRepo(db, config.Cfg)

	// Validate the input data
	if err := userRepo.ValidateSeekerSignUpInput(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email is already taken and whether it's verified or not
	var user models.AuthUser
	err := db.Database(config.Cfg.Cloud.MongoDBName).Collection("auth_users").FindOne(c, bson.M{"email": input.Email}).Decode(&user)

	if err != mongo.ErrNoDocuments {
		// If email exists, check if it's verified
		if user.EmailVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already taken and verified."})
			return
		}

		// If email is not verified, resend the verification email
		token := user.VerificationToken
		verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", config.Cfg.Project.FrontendBaseUrl, token)
		emailBody := fmt.Sprintf(`
			<p>Hello %s,</p>
			<p>Thanks for signing up! Please verify your email by clicking the link below:</p>
			<p><a href="%s">Verify Email</a></p>
			<p>If you did not sign up, you can ignore this email.</p>
		`, input.Email, verificationLink)

		// Send the verification email again
		emailCfg := utils.EmailConfig{
			Host:     config.Cfg.Cloud.EmailHost,
			Port:     config.Cfg.Cloud.EmailPort,
			Username: config.Cfg.Cloud.EmailHostUser,
			Password: config.Cfg.Cloud.EmailHostPassword,
			From:     config.Cfg.Cloud.DefaultFromEmail,
			UseTLS:   config.Cfg.Cloud.EmailUseTLS,
		}

		if err := utils.SendEmail(emailCfg, input.Email, "Verify your email", emailBody); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend verification email", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Please check your email to verify your account. Email resent."})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create the user and send the verification email
	if err := userRepo.CreateSeeker(input, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Seeker", "details": err.Error()})
		return
	}

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
	err := db.(*mongo.Client).Database(config.Cfg.Cloud.MongoDBName).Collection("auth_users").FindOne(c, bson.M{"verification_token": token}).Decode(&user)

	if err == mongo.ErrNoDocuments {
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

	// Update user in the database
	_, err = db.(*mongo.Client).Database(config.Cfg.Cloud.MongoDBName).Collection("auth_users").UpdateOne(
		c,
		bson.M{"_id": user.AuthUserID},
		bson.M{"$set": bson.M{"email_verified": true, "verification_token": ""}},
	)

	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to verify email")
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Email Verified</title>
			<style>
				body { font-family: Arial, sans-serif; background-color: #f2f4f8; color: #333; text-align: center; padding-top: 100px; }
				.card { background: white; padding: 40px; margin: auto; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 90%; max-width: 500px; }
				h1 { color: #28a745; }
				p { margin-top: 10px; font-size: 18px; }
				a { display: inline-block; margin-top: 20px; text-decoration: none; color: white; background-color: #007bff; padding: 10px 20px; border-radius: 5px; }
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
	`))
}

func Login(c *gin.Context) {
	var input dto.LoginInput
	db := c.MustGet("db").(*mongo.Client)
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

	token, err := security.GenerateJWT(user.AuthUserID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
