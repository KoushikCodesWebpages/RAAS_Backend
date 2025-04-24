package repo

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"RAAS/models"
	"gorm.io/gorm"
)


const resetPasskey = "reset@arshan.de" // Change this to your actual passkey

type ResetRequest struct {
	Passkey string `json:"passkey"`
	Email   string `json:"email"` // Added email field to request
}


func ResetDBHandler(c *gin.Context) {
	var req ResetRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.Passkey != resetPasskey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid passkey"})
		return
	}

		// Get auth_user_id by email
		var authUser models.AuthUser
		if err := models.DB.Where("email = ?", req.Email).First(&authUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			}
			return
		}
	
		// Log user id
		log.Printf("🔄 ResetDBHandler triggered for user: %s (ID: %d)", req.Email, authUser.ID)
	

	tables := []string{
		"auth_users",
		"seekers",
		"admins",

		"job_match_scores",
		"user_entry_timelines",
		"selected_job_applications",
		"cover_letters",
		"cv",


	}

	for _, table := range tables {
		if err := models.DB.Exec("DELETE FROM "+table+" WHERE auth_user_id = ?", authUser.ID).Error; err != nil {
			log.Printf("❌ Failed to delete records from table %s: %v", table, err)
		} else {
			log.Printf("✅ Deleted records from table %s for user ID %d", table, authUser.ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ MySQL DB reset and seeded"})
}
