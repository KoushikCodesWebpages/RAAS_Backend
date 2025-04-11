package handlers

import (
	"RAAS/dto"
	"RAAS/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func CreateProfessionalSummary(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.ProfessionalSummaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	skillsJSON, err := json.Marshal(input.Skills)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to process skills",
			"details": err.Error(),
		})
		return
	}

	proSummary := models.ProfessionalSummary{
		AuthUserID:   userID,
		About:        input.About,
		Skills:       datatypes.JSON(skillsJSON),
		AnnualIncome: input.AnnualIncome,
	}

	if err := db.Create(&proSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create professional summary",
			"details": err.Error(),
		})
		return
	}

	// Return the response DTO
	response := dto.ProfessionalSummaryResponse{
		AuthUserID:   proSummary.AuthUserID,
		About:        proSummary.About,
		Skills:       input.Skills,
		AnnualIncome: proSummary.AnnualIncome,
	}

	c.JSON(http.StatusCreated, response)
}

func GetProfessionalSummary(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var proSummary models.ProfessionalSummary
	if err := db.Where("auth_user_id = ?", userID).First(&proSummary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Professional summary not found",
		})
		return
	}

	// Unmarshal skills from JSON to []string
	var skills []string
	if err := json.Unmarshal(proSummary.Skills, &skills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse skills",
			"details": err.Error(),
		})
		return
	}

	response := dto.ProfessionalSummaryResponse{
		AuthUserID:   proSummary.AuthUserID,
		About:        proSummary.About,
		Skills:       skills,
		AnnualIncome: proSummary.AnnualIncome,
	}

	c.JSON(http.StatusOK, response)
}

func UpdateProfessionalSummary(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.ProfessionalSummaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	var proSummary models.ProfessionalSummary
	if err := db.Where("auth_user_id = ?", userID).First(&proSummary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Professional summary not found",
		})
		return
	}

	skillsJSON, err := json.Marshal(input.Skills)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to process skills",
			"details": err.Error(),
		})
		return
	}

	proSummary.About = input.About
	proSummary.Skills = datatypes.JSON(skillsJSON)
	proSummary.AnnualIncome = input.AnnualIncome

	if err := db.Save(&proSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update professional summary",
			"details": err.Error(),
		})
		return
	}

	response := dto.ProfessionalSummaryResponse{
		AuthUserID:   proSummary.AuthUserID,
		About:        proSummary.About,
		Skills:       input.Skills,
		AnnualIncome: proSummary.AnnualIncome,
	}

	c.JSON(http.StatusOK, response)
}
