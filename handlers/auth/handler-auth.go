package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"RAAS/config"
	"RAAS/dto"
	"RAAS/handlers/repo"
	"RAAS/models"
	"RAAS/security"
	"RAAS/utils"
)

func SeekerSignUp(c *gin.Context) {
	var input dto.SeekerSignUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input", "details": err.Error()})
		return
	}

	db := c.MustGet("db").(*mongo.Database)
	userRepo := repo.NewUserRepo(db)

	if err := userRepo.ValidateSeekerSignUpInput(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// Check if email or phone already exists
	emailTaken, phoneTaken, err := userRepo.CheckDuplicateEmailOrPhone(input.Email, input.Number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "check_duplicate_failed", "details": err.Error()})
		return
	}

	if emailTaken || phoneTaken {
		var user models.AuthUser
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.Collection("auth_users").FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email taken", "details": err.Error()})
			return
		}

		if user.EmailVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email_already_verified"})
			return
		}

		// Resend verification email
		verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", config.Cfg.Project.FrontendBaseUrl, user.VerificationToken)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed_to_send_verification_email", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "verification_email_resent"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password_hash_error"})
		return
	}

	// Create the seeker
	if err := userRepo.CreateSeeker(input, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create_seeker_failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "seeker_registered_successfully"})
}

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Missing token")
		return
	}

	db := c.MustGet("db").(*mongo.Database)

	var user models.AuthUser
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.Collection("auth_users").FindOne(ctx, bson.M{"verification_token": token}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.String(http.StatusNotFound, "Invalid or expired token")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	if user.EmailVerified {
		c.String(http.StatusOK, "Email already verified.")
		return
	}

	_, err = db.Collection("auth_users").UpdateOne(
		ctx,
		bson.M{"auth_user_id": user.AuthUserID},
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

	db := c.MustGet("db").(*mongo.Database)
	userRepo := repo.NewUserRepo(db)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := userRepo.AuthenticateUser(ctx, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
		return
	}

	token, err := security.GenerateJWT(user.AuthUserID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token_generation_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
