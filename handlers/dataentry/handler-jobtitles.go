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

	jobTitle := models.PreferredJobTitle{
		PrimaryTitle:   input.PrimaryTitle,
		SecondaryTitle: input.SecondaryTitle,
		TertiaryTitle:  input.TertiaryTitle,
		AuthUserID:     userID,
	}

	if err := h.DB.Create(&jobTitle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job title", "details": err.Error()})
		return
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

	c.JSON(http.StatusOK, jobTitle)
}

// GetJobTitle retrieves the preferred job title for the authenticated user
func (h *JobTitleHandler) GetJobTitle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var jobTitle models.PreferredJobTitle
	if err := h.DB.Where("auth_user_id = ?", userID).First(&jobTitle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No job title found. Please set your preferred job titles first."})
		return
	}

	c.JSON(http.StatusOK, jobTitle)
}

// UpdateJobTitle updates the job titles for the authenticated user
func (h *JobTitleHandler) UpdateJobTitle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var input dto.JobTitleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var jobTitle models.PreferredJobTitle
	if err := h.DB.Where("auth_user_id = ?", userID).First(&jobTitle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job title not found", "details": err.Error()})
		return
	}

	jobTitle.PrimaryTitle = input.PrimaryTitle
	jobTitle.SecondaryTitle = input.SecondaryTitle
	jobTitle.TertiaryTitle = input.TertiaryTitle

	if err := h.DB.Save(&jobTitle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job title", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job title updated successfully", "jobTitle": jobTitle})
}

// PatchJobTitle allows partial update of job titles for the authenticated user
func (h *JobTitleHandler) PatchJobTitle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var input dto.JobTitleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var jobTitle models.PreferredJobTitle
	if err := h.DB.Where("auth_user_id = ?", userID).First(&jobTitle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job title not found", "details": err.Error()})
		return
	}

	if input.PrimaryTitle != "" {
		jobTitle.PrimaryTitle = input.PrimaryTitle
	}
	if input.SecondaryTitle != nil {
		jobTitle.SecondaryTitle = input.SecondaryTitle
	}
	if input.TertiaryTitle != nil {
		jobTitle.TertiaryTitle = input.TertiaryTitle
	}

	if err := h.DB.Save(&jobTitle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to patch job title", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job title patched successfully", "jobTitle": jobTitle})
}
