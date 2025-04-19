package features

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"RAAS/models"
)

// LinkResponseDTO represents the response DTO for job application links
type LinkResponseDTO struct {
	JobID   string `json:"job_id"`
	JobLink string `json:"job_link"`
	Source  string `json:"source"`
}

// LinkProviderHandler handles requests for job application links
type LinkProviderHandler struct {
	db *gorm.DB
}


// NewLinkProviderHandler returns a new instance of LinkProviderHandler
func NewLinkProviderHandler(db *gorm.DB) *LinkProviderHandler {
	return &LinkProviderHandler{db: db}
}

// PostAndGetLink handles POST requests to retrieve job application links
func (h *LinkProviderHandler) PostAndGetLink(c *gin.Context) {
	var req struct {
		JobID string `json:"job_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Failed to bind request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid job_id in request body"})
		return
	}
	jobID := req.JobID

	authUserID, ok := c.MustGet("userID").(uuid.UUID)
	if !ok {
		fmt.Println("❌ Failed to get userID from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID from JWT claims"})
		return
	}

	// Check if the job was selected by the user
	var selectedJob models.SelectedJobApplication
	if err := h.db.Where("auth_user_id = ? AND job_id = ?", authUserID, jobID).First(&selectedJob).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("🚫 Job not selected by the user")
			c.JSON(http.StatusForbidden, gin.H{"error": "Job not selected by the user"})
			return
		}
		fmt.Println("❌ DB error while checking selected job:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify job selection"})
		return
	}

	// Retrieve JobLink from the unified Job model
	var job models.Job
	if err := h.db.Where("job_id = ?", jobID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("🚫 Job not found in jobs table")
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		fmt.Println("❌ DB error while fetching job:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job info"})
		return
	}

	// Build and return the response
	response := LinkResponseDTO{
		JobID:   job.JobID,
		JobLink: job.JobLink,
		Source:  job.Source,
	}

	// Update view_link = true
	if err := h.db.Model(&models.SelectedJobApplication{}).
		Where("auth_user_id = ? AND job_id = ?", authUserID, jobID).
		Update("view_link", true).Error; err != nil {
		fmt.Println("⚠️ Failed to update view_link field:", err)
	}

	c.JSON(http.StatusOK, response)
}
