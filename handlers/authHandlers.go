package handlers

import (

    "RAAS/dto"
    "RAAS/security"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "net/http"
    "RAAS/config"
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


// func GoogleLoginHandler(c *gin.Context) {
// 	url := security.GoogleOAuthConfig.AuthCodeURL("random-state-string") // you can use a CSRF-safe state later
// 	c.Redirect(http.StatusTemporaryRedirect, url)
// }

// func GoogleCallbackHandler(c *gin.Context) {
//     code := c.Query("code")

//     token, err := security.GoogleOAuthConfig.Exchange(context.Background(), code)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
//         return
//     }

//     client := security.GoogleOAuthConfig.Client(context.Background(), token)
//     resp, _ := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
//     defer resp.Body.Close()

//     body, _ := io.ReadAll(resp.Body)
//     var userInfo security.UserInfo
//     json.Unmarshal(body, &userInfo)

//     user, err := createSeekerFromGoogleOAuth(db, userInfo)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
//         return
//     }

//     tokenString := security.GenerateJWT(user.ID)

//     // You could redirect to frontend with token
//     c.JSON(http.StatusOK, gin.H{
//         "token": tokenString,
//         "user":  userInfo,
//     })
// }


