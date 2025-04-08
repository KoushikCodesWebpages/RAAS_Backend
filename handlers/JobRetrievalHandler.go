package handlers

import (
	"RAAS/models"
	"RAAS/dto"
	"strings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// JobRetrievalHandler handles GET /api/jobs — returns all Xing and LinkedIn jobs, optionally filtered by title
func JobRetrievalHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	titleFilter := strings.ToLower(c.Query("title"))

	var linkedinJobs []models.LinkedInJobMetaData
	var xingJobs []models.XingJobMetaData

	// Apply filter if provided
	if titleFilter != "" {
		db.Where("LOWER(title) LIKE ?", "%"+titleFilter+"%").Find(&linkedinJobs)
		db.Where("LOWER(title) LIKE ?", "%"+titleFilter+"%").Find(&xingJobs)
	} else {
		db.Find(&linkedinJobs)
		db.Find(&xingJobs)
	}

	// Merge into DTOs
	var jobs []dto.JobDTO

	for _, job := range linkedinJobs {
		jobs = append(jobs, dto.JobDTO{
			Source:     "linkedin",
			ID:         job.ID,
			JobID:      job.JobID,
			Title:      job.Title,
			Company:    job.Company,
			Location:   job.Location,
			PostedDate: job.PostedDate,
			Link:       job.Link,
			Processed:  job.Processed,
		})
	}

	for _, job := range xingJobs {
		jobs = append(jobs, dto.JobDTO{
			Source:     "xing",
			ID:         job.ID,
			JobID:      job.JobID,
			Title:      job.Title,
			Company:    job.Company,
			Location:   job.Location,
			PostedDate: job.PostedDate,
			Link:       job.Link,
			Processed:  job.Processed,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
		"pagination": gin.H{ // Placeholder for future pagination support
			"total": len(jobs),
		},
	})
}
