package features

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"fmt" // Add fmt package for logging
)

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

	// Debugging: Log the preferred job titles
	fmt.Println("Preferred Job Titles:", preferred)

	// Collect preferred titles
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

	// Debugging: Log the collected preferred titles
	fmt.Println("Collected Preferred Titles:", preferredTitles)

	if len(preferredTitles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No preferred job titles set for user."})
		return
	}

	// Title filtering (case-insensitive)
	var conditions []string
	var values []interface{}
	for _, title := range preferredTitles {
		conditions = append(conditions, "LOWER(title) LIKE ?")
		values = append(values, "%"+strings.ToLower(title)+"%")
	}
	whereClause := strings.Join(conditions, " OR ")

	// Debugging: Log the generated SQL conditions and values
	fmt.Println("SQL WHERE Clause:", whereClause)
	fmt.Println("SQL Values:", values)

	// Query jobs from both sources
	var linkedinJobs []models.LinkedInJobMetaData
	var xingJobs []models.XingJobMetaData

	// Log the queries being executed
	db.Debug().Where(whereClause, values...).Find(&linkedinJobs) // Debugging SQL query
	db.Debug().Where(whereClause, values...).Find(&xingJobs) // Debugging SQL query

	// Log the number of jobs fetched
	fmt.Println("Number of LinkedIn Jobs Found:", len(linkedinJobs))
	fmt.Println("Number of Xing Jobs Found:", len(xingJobs))

	// LinkedIn Jobs
	for _, job := range linkedinJobs {
		var jobDesc models.LinkedInJobDescription
		if err := db.Where("job_id = ?", job.JobID).First(&jobDesc).Error; err != nil {
			// Debugging: Log the error when fetching job descriptions
			fmt.Println("Error fetching LinkedIn job description for job_id:", job.JobID, "Error:", err)
			continue
		}

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
		})
	}

	// Xing Jobs
	for _, job := range xingJobs {
		var jobDesc models.XingJobDescription
		if err := db.Where("job_id = ?", job.JobID).First(&jobDesc).Error; err != nil {
			// Debugging: Log the error when fetching job descriptions
			fmt.Println("Error fetching Xing job description for job_id:", job.JobID, "Error:", err)
			continue
		}

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
		})
	}

	// Debugging: Log the total number of jobs collected
	fmt.Println("Total Jobs Retrieved:", len(jobs))

	// Respond with the jobs
	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
		"pagination": gin.H{
			"total": len(jobs),
		},
	})
}
