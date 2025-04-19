package dataentry

import (
	"RAAS/models"
	"RAAS/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"github.com/google/uuid"
)

// JobTitleHandler struct
type JobTitleHandler struct {
	DB *gorm.DB
}

// NewJobTitleHandler creates a new JobTitleHandler
func NewJobTitleHandler(db *gorm.DB) *JobTitleHandler {
	return &JobTitleHandler{DB: db}
}

// CreateJobTitle creates a new preferred job title for the authenticated user
func (h *JobTitleHandler) CreateJobTitle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.JobTitleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find the Seeker model by AuthUserID
	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Set the job titles for the authenticated user
	seeker.PrimaryTitle = input.PrimaryTitle
	seeker.SecondaryTitle = input.SecondaryTitle
	seeker.TertiaryTitle = input.TertiaryTitle

	// Save the updated Seeker model
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job title", "details": err.Error()})
		return
	}

	// Create the response using JobTitleResponse
	jobTitleResponse := dto.JobTitleResponse{
		AuthUserID:   seeker.AuthUserID,
		PrimaryTitle: seeker.PrimaryTitle,
		SecondaryTitle: seeker.SecondaryTitle,
		TertiaryTitle: seeker.TertiaryTitle,
	}

	// Update the UserEntryTimeline - set PreferredJobTitlesCompleted to true
	var timeline models.UserEntryTimeline
	if err := h.DB.First(&timeline, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
		return
	}

	timeline.PreferredJobTitlesCompleted = true

	// Save the updated timeline
	if err := h.DB.Save(&timeline).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timeline", "details": err.Error()})
		return
	}

	// Return the formatted response
	c.JSON(http.StatusOK, jobTitleResponse)
}


// GetJobTitle retrieves the preferred job title for the authenticated user
func (h *JobTitleHandler) GetJobTitle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var seeker models.Seeker
	if err := h.DB.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No job title found. Please set your preferred job titles first."})
		return
	}

	// Create response containing the job titles from Seeker
	jobTitleResponse := dto.JobTitleResponse{
		AuthUserID:  seeker.AuthUserID,
		PrimaryTitle: seeker.PrimaryTitle,
		SecondaryTitle: seeker.SecondaryTitle,
		TertiaryTitle:  seeker.TertiaryTitle,
	}

	c.JSON(http.StatusOK, jobTitleResponse)
}

// UpdateJobTitle updates the job titles for the authenticated user
func (h *JobTitleHandler) UpdateJobTitle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var input dto.JobTitleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find the Seeker model by AuthUserID
	var seeker models.Seeker
	if err := h.DB.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found", "details": err.Error()})
		return
	}

	// Update the job titles for the Seeker
	seeker.PrimaryTitle = input.PrimaryTitle
	seeker.SecondaryTitle = input.SecondaryTitle
	seeker.TertiaryTitle = input.TertiaryTitle

	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job title", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job title updated successfully", "jobTitle": seeker})
}

// PatchJobTitle allows partial update of job titles for the authenticated user
func (h *JobTitleHandler) PatchJobTitle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var input dto.JobTitleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find the Seeker model by AuthUserID
	var seeker models.Seeker
	if err := h.DB.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found", "details": err.Error()})
		return
	}

	// Update only the fields that are provided
	if input.PrimaryTitle != "" {
		seeker.PrimaryTitle = input.PrimaryTitle
	}
	if input.SecondaryTitle != nil {
		seeker.SecondaryTitle = input.SecondaryTitle
	}
	if input.TertiaryTitle != nil {
		seeker.TertiaryTitle = input.TertiaryTitle
	}

	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to patch job title", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job title patched successfully", "jobTitle": seeker})
}
