package auth

// import (
// 	"net/http"
// 	"RAAS/models"
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// // VerifyEmailHandler handles email verification when user clicks the link
// func VerifyEmailHandler(c *gin.Context) {
// 	db := c.MustGet("db").(*gorm.DB)
// 	token := c.Query("token")

// 	if token == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
// 		return
// 	}

// 	var user models.AuthUser
// 	if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
// 		return
// 	}

// 	user.EmailVerified = true
// 	user.VerificationToken = "" // Invalidate the token

// 	if err := db.Save(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user verification status"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully!"})
// }
