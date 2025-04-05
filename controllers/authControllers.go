package controllers

import (
    "RAAS/repositories"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "net/http"
    "fmt"
    "RAAS/security"
)

// SeekerSignUp handles the Seeker registration process
func SeekerSignUp(c *gin.Context, db *gorm.DB) {
    var input repositories.SeekerSignUpInput

    // Bind JSON body to the input struct
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    // Validate input
    if err := repositories.ValidateSeekerSignUpInput(input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if email is already registered
    emailTaken, err := repositories.IsEmailTaken(db, input.Email)
    if err != nil {
        fmt.Println("Database error in IsEmailTaken:", err) // 👈 Debug print
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
        return
    }
    if emailTaken {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already taken"})
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        fmt.Println("Error hashing password:", err) // 👈 Debug print
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    // Create the user
    err = repositories.CreateSeeker(db, input, string(hashedPassword))
    if err != nil {
        fmt.Println("Error creating seeker:", err) // 👈 Debug print
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Seeker", "details": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Seeker registered successfully"})

}



type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// Login handles user authentication
func Login(c *gin.Context) {
    var input LoginInput
    db := c.MustGet("db").(*gorm.DB) // Ensure DB instance is available

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    user, err := repositories.AuthenticateUser(db, input.Email, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Generate JWT
    token, err := security.GenerateJWT(user.Email, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token, "role": user.Role})
}