package handlers

import (

    "RAAS/dto"
    "RAAS/security"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "net/http"
    "RAAS/config"
    
)

func SeekerSignUp(c *gin.Context, db *gorm.DB) {
    var input dto.SeekerSignUpInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    if err := validateSeekerSignUpInput(input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

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

    // ✅ Now pass cfg to createSeeker
    if err := createSeeker(db, input, string(hashedPassword), cfg); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Seeker", "details": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Seeker registered successfully. Please check your email to verify."})
}



func Login(c *gin.Context) {
    var input dto.LoginInput
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

    token, err := security.GenerateJWT(user.Email, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token, "role": user.Role})
}


