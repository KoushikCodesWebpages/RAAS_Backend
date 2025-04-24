package dataentry

import (
	"RAAS/dto"
	"RAAS/handlers"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProfessionalSummaryHandler struct
type ProfessionalSummaryHandler struct {
	DB *gorm.DB
}

// NewProfessionalSummaryHandler creates a new handler instance
func NewProfessionalSummaryHandler(db *gorm.DB) *ProfessionalSummaryHandler {
	return &ProfessionalSummaryHandler{DB: db}
}






func (h *ProfessionalSummaryHandler) CreateProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.ProfessionalSummaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find Seeker for the given user using the utility function
	seeker, err := handlers.FindSeekerByUserID(h.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Check if ProfessionalSummary is already filled using the utility function
	isFilled, err := handlers.IsFieldFilled(seeker.ProfessionalSummary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check professional summary", "details": err.Error()})
		return
	}

	if isFilled {
		c.JSON(http.StatusConflict, gin.H{"error": "Professional summary already filled"})
		return
	}

	// Prepare the professional summary map
	professionalSummary := &dto.ProfessionalSummaryRequest{
		About:        input.About,
		Skills:       input.Skills,
		AnnualIncome: input.AnnualIncome,
	}

	// Set the professional summary using the utility function
	if err := handlers.SetProfessionalSummary(seeker, professionalSummary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set professional summary", "details": err.Error()})
		return
	}

	// Save the updated Seeker record
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update professional summary", "details": err.Error()})
		return
	}

	// Return the ProfessionalSummaryResponse
	professionalSummaryResponse := dto.ProfessionalSummaryResponse{
		AuthUserID:   seeker.AuthUserID,
		About:        input.About,
		Skills:       input.Skills,
		AnnualIncome: input.AnnualIncome,
	}

	c.JSON(http.StatusCreated, professionalSummaryResponse)
}





func (h *ProfessionalSummaryHandler) GetProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	seeker, err := handlers.FindSeekerByUserID(h.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	isFilled, err := handlers.IsFieldFilled(seeker.ProfessionalSummary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check professional summary", "details": err.Error()})
		return
	}
	if !isFilled {
		c.JSON(http.StatusNotFound, gin.H{"error": "Professional summary not yet filled"})
		return
	}

	profSummary, err := handlers.GetProfessionalSummary(seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse professional summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ProfessionalSummaryResponse{
		AuthUserID:   seeker.AuthUserID,
		About:        profSummary.About,
		Skills:       profSummary.Skills,
		AnnualIncome: profSummary.AnnualIncome,
	})
}







func (h *ProfessionalSummaryHandler) UpdateProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.ProfessionalSummaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	seeker, err := handlers.FindSeekerByUserID(h.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	if err := handlers.SetProfessionalSummary(seeker, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update professional summary", "details": err.Error()})
		return
	}

	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database save failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ProfessionalSummaryResponse{
		AuthUserID:   seeker.AuthUserID,
		About:        input.About,
		Skills:       input.Skills,
		AnnualIncome: input.AnnualIncome,
	})
}




func (h *ProfessionalSummaryHandler) PatchProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	seeker, err := handlers.FindSeekerByUserID(h.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	var professionalSummary map[string]interface{}
	if err := json.Unmarshal(seeker.ProfessionalSummary, &professionalSummary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal professional summary", "details": err.Error()})
		return
	}

	for key, value := range updates {
		professionalSummary[key] = value
	}

	updatedJSON, err := json.Marshal(professionalSummary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated professional summary", "details": err.Error()})
		return
	}

	seeker.ProfessionalSummary = updatedJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated professional summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Professional summary updated"})
}

