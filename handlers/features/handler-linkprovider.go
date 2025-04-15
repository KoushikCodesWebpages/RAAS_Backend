package features

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"fmt"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid job_id in request body"})
		return
	}
	jobID := req.JobID

	// Extract user information from JWT claims
	authUserID, ok := c.MustGet("userID").(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID from JWT claims"})
		return
	}

	fmt.Println("Auth User ID:", authUserID)
	fmt.Println("Job ID:", jobID)

	// Verify if the job is selected by the user
	var selectedJobApplication models.SelectedJobApplication
	result := h.db.Where("auth_user_id = ? AND job_id = ?", authUserID, jobID).First(&selectedJobApplication)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Job not selected by the user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify job selection"})
		return
	}

	fmt.Println("Selected Job Application:", selectedJobApplication)

	// Retrieve LinkedIn job application link
	var linkedInLink models.LinkedInJobApplicationLink
	result = h.db.Where("job_id = ?", jobID).First(&linkedInLink)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve LinkedIn job application link"})
		return
	}

	fmt.Println("LinkedIn Link:", linkedInLink)

	// Retrieve Xing job application link
	var xingLink models.XingJobApplicationLink
	result = h.db.Where("job_id = ?", jobID).First(&xingLink)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Xing job application link"})
		return
	}

	fmt.Println("Xing Link:", xingLink)

	// Prepare response
	var response []LinkResponseDTO
	if linkedInLink.JobID != "" {
		response = append(response, LinkResponseDTO{
			JobID:   linkedInLink.JobID,
			JobLink: linkedInLink.JobLink,
			Source:  "LinkedIn",
		})
	}
	if xingLink.JobID != "" {
		response = append(response, LinkResponseDTO{
			JobID:   xingLink.JobID,
			JobLink: xingLink.JobLink,
			Source:  "Xing",
		})
	}

	fmt.Println("Response:", response)

	// Return response
	if len(response) > 0 {
		h.db.Model(&models.SelectedJobApplication{}).Where("auth_user_id = ? AND job_id = ?", authUserID, jobID).Update("view_link", true)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "No job application links found"})
	}
}