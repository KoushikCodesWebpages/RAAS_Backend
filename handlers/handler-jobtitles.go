package handlers

import (
	"RAAS/models"
	"RAAS/dto"
	//"RAAS/security"
	//"RAAS/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"github.com/google/uuid"
)

// CreateJobTitle creates a new preferred job title for the authenticated user
func CreateJobTitle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input dto.JobTitleInput
	userID := c.MustGet("userID").(uuid.UUID) // From JWT claims

	// Bind the input JSON to the JobTitleInput struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Create PreferredJobTitle
	jobTitle := models.PreferredJobTitle{
		PrimaryTitle:   input.PrimaryTitle,
		SecondaryTitle: input.SecondaryTitle,
		TertiaryTitle:  input.TertiaryTitle,
		AuthUserID:     userID,
	}

	if err := db.Create(&jobTitle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job title", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobTitle)

}

// GetJobTitles retrieves the preferred job titles for the authenticated user
func GetJobTitle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID) // From JWT

	var jobTitle models.PreferredJobTitle
	if err := db.Where("auth_user_id = ?", userID).First(&jobTitle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No job title found. Please set your preferred job titles first."})
		return
	}

	c.JSON(http.StatusOK, jobTitle)
}

// UpdateJobTitle updates the job titles for the authenticated user
func UpdateJobTitle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input dto.JobTitleInput
	userID := c.MustGet("userID").(uuid.UUID) // From JWT claims

	// Bind the input JSON to the JobTitleInput struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find the existing job title record (only one, as we assume one record per user)
	var jobTitle models.PreferredJobTitle
	if err := db.Where("auth_user_id = ?", userID).First(&jobTitle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job title not found", "details": err.Error()})
		return
	}

	// Update the job title fields
	jobTitle.PrimaryTitle = input.PrimaryTitle
	jobTitle.SecondaryTitle = input.SecondaryTitle
	jobTitle.TertiaryTitle = input.TertiaryTitle

	if err := db.Save(&jobTitle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job title", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job title updated successfully", "jobTitle": jobTitle})
}

// PatchJobTitle allows partial update of job titles for the authenticated user
func PatchJobTitle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input dto.JobTitleInput
	userID := c.MustGet("userID").(uuid.UUID) // From JWT claims

	// Bind the input JSON to the JobTitleInput struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find the existing job title record (only one, as we assume one record per user)
	var jobTitle models.PreferredJobTitle
	if err := db.Where("auth_user_id = ?", userID).First(&jobTitle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job title not found", "details": err.Error()})
		return
	}

	// Update only the fields that are provided
	if input.PrimaryTitle != "" {
		jobTitle.PrimaryTitle = input.PrimaryTitle
	}
	if input.SecondaryTitle != nil {
		jobTitle.SecondaryTitle = input.SecondaryTitle
	}
	if input.TertiaryTitle != nil {
		jobTitle.TertiaryTitle = input.TertiaryTitle
	}

	if err := db.Save(&jobTitle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to patch job title", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job title patched successfully", "jobTitle": jobTitle})
}
