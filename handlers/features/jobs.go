package features

import (
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobDataHandler struct {
	db *gorm.DB
}

func NewJobDataHandler(db *gorm.DB) *JobDataHandler {
	return &JobDataHandler{db: db}
}

// Get all jobs (combined from LinkedIn, Xing, and other sources)
// Get all jobs (combined from LinkedIn, Xing, and other sources)
func (h *JobDataHandler) GetAllJobs(c *gin.Context) {
	var jobs []models.Job

	// Fetch all jobs from the database
	if err := h.db.Find(&jobs).Error; err != nil {
		// Return an internal server error if fetching the jobs fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	// If no jobs are found
	if len(jobs) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No jobs available"})
		return
	}

	// Return all jobs
	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}
