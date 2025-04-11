package handlers

import (
	"RAAS/models"
	"RAAS/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"github.com/google/uuid"
)

func CreatePersonalInfo(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.PersonalInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	personalInfo := models.PersonalInfo{
		AuthUserID:      userID,
		FirstName:       input.FirstName,
		SecondName:      input.SecondName,
		DateOfBirth:     input.DateOfBirth,
		Address:         input.Address,
		LinkedInProfile: input.LinkedInProfile,
	}

	if err := db.Create(&personalInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create personal info",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, personalInfo)
}

func GetPersonalInfo(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var personalInfo models.PersonalInfo
	if err := db.Where("auth_user_id = ?", userID).First(&personalInfo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal info not found"})
		return
	}

	c.JSON(http.StatusOK, personalInfo)
}

func UpdatePersonalInfo(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.PersonalInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var personalInfo models.PersonalInfo
	if err := db.Where("auth_user_id = ?", userID).First(&personalInfo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal info not found"})
		return
	}

	personalInfo.FirstName = input.FirstName
	personalInfo.SecondName = input.SecondName
	personalInfo.DateOfBirth = input.DateOfBirth
	personalInfo.Address = input.Address
	personalInfo.LinkedInProfile = input.LinkedInProfile

	if err := db.Save(&personalInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, personalInfo)
}

func PatchPersonalInfo(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Only allow specific fields
	allowedFields := map[string]bool{
		"firstName":       true,
		"secondName":      true,
		"dob":             true,
		"address":         true,
		"linkedinProfile": true,
	}

	for k := range updates {
		if !allowedFields[k] {
			delete(updates, k)
		}
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	if err := db.Model(&models.PersonalInfo{}).
		Where("auth_user_id = ?", userID).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to patch personal info", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Personal info updated"})
}
