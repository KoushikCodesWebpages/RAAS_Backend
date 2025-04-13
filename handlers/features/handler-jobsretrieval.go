package features

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"
	"strings"
	//"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JobRetrievalHandler retrieves jobs based on the user's preferred job titles
func JobRetrievalHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var jobs []dto.JobDTO

	// Fetch user's preferred job titles
	var preferred models.PreferredJobTitle
	if err := db.Where("auth_user_id = ?", userID).First(&preferred).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Set Your Job Title First",
		})
		return
	}

	// Collect preferred titles (non-nil, non-empty)
	var preferredTitles []string
	if preferred.PrimaryTitle != "" {
		preferredTitles = append(preferredTitles, preferred.PrimaryTitle)
	}
	if preferred.SecondaryTitle != nil && *preferred.SecondaryTitle != "" {
		preferredTitles = append(preferredTitles, *preferred.SecondaryTitle)
	}
	if preferred.TertiaryTitle != nil && *preferred.TertiaryTitle != "" {
		preferredTitles = append(preferredTitles, *preferred.TertiaryTitle)
	}

	if len(preferredTitles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No preferred job titles set for user."})
		return
	}

	// Build WHERE clause for case-insensitive title matches
	var conditions []string
	var values []interface{}
	for _, title := range preferredTitles {
		conditions = append(conditions, "LOWER(title) LIKE ?")
		values = append(values, "%"+strings.ToLower(title)+"%")
	}
	whereClause := strings.Join(conditions, " OR ")

	// Query jobs from both sources
	var linkedinJobs []models.LinkedInJobMetaData
	var xingJobs []models.XingJobMetaData

	db.Where(whereClause, values...).Find(&linkedinJobs)
	db.Where(whereClause, values...).Find(&xingJobs)

	// LinkedIn Jobs
	for _, job := range linkedinJobs {
		var jobDesc models.LinkedInJobDescription
		if err := db.Where("job_id = ?", job.JobID).First(&jobDesc).Error; err != nil {
			continue
		}

		matchScore := 75.0

		jobs = append(jobs, dto.JobDTO{
			Source:     "linkedin",
			ID:         job.ID,
			JobID:      job.JobID,
			Title:      job.Title,
			Company:    job.Company,
			Location:   job.Location,
			PostedDate: job.PostedDate,
			Processed:  job.Processed,
			JobType:    jobDesc.JobType,
			Skills:     jobDesc.Skills,
			MatchScore: matchScore,
		})
	}

	// Xing Jobs
	for _, job := range xingJobs {
		var jobDesc models.XingJobDescription
		if err := db.Where("job_id = ?", job.JobID).First(&jobDesc).Error; err != nil {
			continue
		}

		matchScore := 75.0

		jobs = append(jobs, dto.JobDTO{
			Source:     "xing",
			ID:         job.ID,
			JobID:      job.JobID,
			Title:      job.Title,
			Company:    job.Company,
			Location:   job.Location,
			PostedDate: job.PostedDate,
			Processed:  job.Processed,
			JobType:    jobDesc.JobType,
			Skills:     jobDesc.Skills,
			MatchScore: matchScore,
		})
	}

	// Send back jobs and their match score
	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
		"pagination": gin.H{
			"total": len(jobs),
		},
	})
}
