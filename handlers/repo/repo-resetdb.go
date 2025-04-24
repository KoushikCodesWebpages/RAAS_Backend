package repo

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"RAAS/models"
)

const resetPasskey = "reset@arshan.de"

type ResetRequest struct {
	Passkey string `json:"passkey"`
	Email   string `json:"email"`
}

func ResetDBHandler(c *gin.Context) {
	var req ResetRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil || req.Passkey != resetPasskey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid passkey or bad request"})
		return
	}

	// Fetch user by email
	var authUser models.AuthUser
	if err := models.DB.Where("email = ?", req.Email).First(&authUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Printf("❌ DB error retrieving user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	userID := authUser.ID
	log.Printf("🔄 Reset triggered for user: %s (ID: %s)", req.Email, userID)

	if err := models.DB.Unscoped().Where("email = ?", req.Email).Delete(&models.AuthUser{}).Error; err != nil {
		log.Printf("❌ Failed to delete user from auth_users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Tables to check
	tables := []string{
		"seekers",
		"admins",
		"match_scores",
		"user_entry_timelines",
		"selected_job_applications",
		"cover_letters",
		"cv",
	}

	leftovers := []string{}

	// Clean each table
	for _, table := range tables {
		var count int64

		// Count records before deletion
		if err := models.DB.Table(table).Where("auth_user_id = ?", userID).Count(&count).Error; err != nil {
			log.Printf("⚠️ Error checking %s: %v", table, err)
			continue
		}

		if count > 0 {
			// Attempt deletion
			if err := models.DB.Exec("DELETE FROM "+table+" WHERE auth_user_id = ?", userID).Error; err != nil {
				log.Printf("❌ Error deleting from %s: %v", table, err)
			} else {
				log.Printf("✅ Deleted %d rows from %s", count, table)
			}
		}

		// Re-check for leftovers
		models.DB.Table(table).Where("auth_user_id = ?", userID).Count(&count)
		if count > 0 {
			leftovers = append(leftovers, table)
		}
	}

	// Final response
	if len(leftovers) > 0 {
		log.Printf("⚠️ Leftover data found in: %v", leftovers)
		c.JSON(http.StatusOK, gin.H{
			"message":         "Partially deleted user data. Leftovers detected.",
			"leftover_tables": leftovers,
		})
	} else {
		log.Println("✅ All user data removed successfully.")
		c.JSON(http.StatusOK, gin.H{
			"message": "✅ User and all associated data deleted successfully.",
		})
	}
}
