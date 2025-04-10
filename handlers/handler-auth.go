package handlers

import (

    "RAAS/dto"
    "RAAS/security"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "net/http"
    "RAAS/config"

	"RAAS/models"
    // "context"
    // "encoding/json"
    // "io"
    
)

func SeekerSignUp(c *gin.Context) {
    var input dto.SeekerSignUpInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    if err := validateSeekerSignUpInput(input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // ✅ Get DB from Gin context
    db := c.MustGet("db").(*gorm.DB)

    emailTaken, err := isEmailTaken(db, input.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
        return
    }
    if emailTaken {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already taken"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    cfg, err := config.InitConfig() // 👈 Load config
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load config", "details": err.Error()})
        return
    }

    if err := createSeeker(db, input, string(hashedPassword), cfg); err != nil {
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
	if err := db.(*gorm.DB).Where("verification_token = ?", token).First(&user).Error; err != nil {
		c.String(http.StatusNotFound, "Invalid or expired token")
		return
	}

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
    cfg, _ := config.InitConfig()
    db := c.MustGet("db").(*gorm.DB)

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    user, err := authenticateUser(db, input.Email, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    token, err := security.GenerateJWT(user.ID, user.Email, user.Role, cfg)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

