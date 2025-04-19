package repo

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"RAAS/models"
)

const resetPasskey = "reset@arshan.de" // Change this to your actual passkey

type ResetRequest struct {
	Passkey string `json:"passkey"`
}

func ResetDBHandler(c *gin.Context) {
	var req ResetRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.Passkey != resetPasskey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid passkey"})
		return
	}

	log.Println("🔄 ResetDBHandler triggered with valid passkey...")

	tables := []string{
		"auth_users",
		"seekers",
		"admins",

		"personal_infos",
		"professional_summaries",
		"work_experiences",
		"educations",
		"languages",
		"certificates",
		"preferred_job_titles",

		"job_match_scores",
	}

	models.ResetDB(models.DB, tables)
	models.AutoMigrate()
	models.SeedJobs(models.DB)

	c.JSON(http.StatusOK, gin.H{"message": "✅ MySQL DB reset and seeded"})
}
