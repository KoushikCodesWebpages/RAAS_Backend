package dataentry

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

// ProfessionalSummaryHandler struct
type ProfessionalSummaryHandler struct {
	DB *gorm.DB
}

// NewProfessionalSummaryHandler creates a new ProfessionalSummaryHandler
func NewProfessionalSummaryHandler(db *gorm.DB) *ProfessionalSummaryHandler {
	return &ProfessionalSummaryHandler{DB: db}
}

// CreateProfessionalSummary creates a professional summary for the authenticated user
func (h *ProfessionalSummaryHandler) CreateProfessionalSummary(c *gin.Context) {
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

	if err := h.DB.Create(&proSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create professional summary",
			"details": err.Error(),
		})
		return
	}

	var timeline models.UserEntryTimeline
    if err := h.DB.First(&timeline, "user_id = ?", userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
        return
    }

    timeline.ProfessionalSummariesCompleted = true

    // Save the updated timeline
    if err := h.DB.Save(&timeline).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timeline", "details": err.Error()})
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

// GetProfessionalSummary retrieves the professional summary of the authenticated user
func (h *ProfessionalSummaryHandler) GetProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var proSummary models.ProfessionalSummary
	if err := h.DB.Where("auth_user_id = ?", userID).First(&proSummary).Error; err != nil {
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

// UpdateProfessionalSummary updates the professional summary of the authenticated user
func (h *ProfessionalSummaryHandler) UpdateProfessionalSummary(c *gin.Context) {
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
	if err := h.DB.Where("auth_user_id = ?", userID).First(&proSummary).Error; err != nil {
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

	if err := h.DB.Save(&proSummary).Error; err != nil {
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
