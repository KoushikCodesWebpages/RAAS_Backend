package dataentry

import (
	"RAAS/models"
	"RAAS/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"net/http"
	"encoding/json"
)

// ProfessionalSummaryHandler struct
type ProfessionalSummaryHandler struct {
	DB *gorm.DB
}

// NewProfessionalSummaryHandler creates a new handler instance
func NewProfessionalSummaryHandler(db *gorm.DB) *ProfessionalSummaryHandler {
	return &ProfessionalSummaryHandler{DB: db}
}

// CreateProfessionalSummary handles the creation of professional summary within Seeker model
func (h *ProfessionalSummaryHandler) CreateProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.ProfessionalSummaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find Seeker for the given user
	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Check if ProfessionalSummary is already filled
	if len(seeker.ProfessionalSummary) > 0 {
		var existingSummary map[string]interface{}
		if err := json.Unmarshal(seeker.ProfessionalSummary, &existingSummary); err == nil && len(existingSummary) > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Professional summary already filled"})
			return
		}
	}

	// Prepare the professional summary map
	professionalSummary := map[string]interface{}{
		"about":         input.About,
		"skills":        input.Skills,
		"annualIncome":  input.AnnualIncome,
	}

	// Marshal the professional summary into JSON
	professionalSummaryJSON, err := json.Marshal(professionalSummary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal professional summary", "details": err.Error()})
		return
	}

	// Update Seeker's ProfessionalSummary field
	seeker.ProfessionalSummary = professionalSummaryJSON

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

// GetProfessionalSummary retrieves the professional summary of the authenticated user
func (h *ProfessionalSummaryHandler) GetProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Check if ProfessionalSummary is empty or uninitialized
	if len(seeker.ProfessionalSummary) == 0 || string(seeker.ProfessionalSummary) == "null" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Professional summary not yet filled"})
		return
	}

	// Unmarshal the ProfessionalSummary JSON field into the struct
	var professionalSummary map[string]interface{}
	if err := json.Unmarshal(seeker.ProfessionalSummary, &professionalSummary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse professional summary", "details": err.Error()})
		return
	}

	// Ensure all required fields exist and are of the correct type
	about, aboutOk := professionalSummary["about"].(string)
	skills, skillsOk := professionalSummary["skills"].([]interface{})
	annualIncome, annualIncomeOk := professionalSummary["annualIncome"].(float64)

	if !aboutOk || !skillsOk || !annualIncomeOk {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid professional summary data",
			"details": "Expected 'about' to be a string, 'skills' to be an array, and 'annualIncome' to be a float64",
		})
		return
	}

	// Convert skills from []interface{} to []string
	skillsStr := make([]string, len(skills))
	for i, skill := range skills {
		if skillStr, ok := skill.(string); ok {
			skillsStr[i] = skillStr
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid skill data type"})
			return
		}
	}

	// Create a response containing the ProfessionalSummary data
	professionalSummaryResponse := dto.ProfessionalSummaryResponse{
		AuthUserID:   seeker.AuthUserID,
		About:        about,
		Skills:       skillsStr,
		AnnualIncome: annualIncome,
	}

	c.JSON(http.StatusOK, professionalSummaryResponse)
}


// UpdateProfessionalSummary performs a full update on the professional summary within Seeker model
func (h *ProfessionalSummaryHandler) UpdateProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.ProfessionalSummaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Prepare the professional summary map
	professionalSummary := map[string]interface{}{
		"about":         input.About,
		"skills":        input.Skills,
		"annualIncome":  input.AnnualIncome,
	}

	// Marshal the professional summary into JSON
	professionalSummaryJSON, err := json.Marshal(professionalSummary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal professional summary", "details": err.Error()})
		return
	}

	// Update Seeker's ProfessionalSummary field
	seeker.ProfessionalSummary = professionalSummaryJSON

	// Save the updated Seeker record
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update professional summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seeker)
}

// PatchProfessionalSummary handles partial updates to professional summary within Seeker model
func (h *ProfessionalSummaryHandler) PatchProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Convert the ProfessionalSummary field back to a map (deserialize)
	var professionalSummary map[string]interface{}
	if err := json.Unmarshal(seeker.ProfessionalSummary, &professionalSummary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal professional summary", "details": err.Error()})
		return
	}

	// Update the map with the valid updates
	for key, value := range updates {
		professionalSummary[key] = value
	}

	// Marshal the updated professional summary back to JSON
	updatedProfessionalSummaryJSON, err := json.Marshal(professionalSummary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated professional summary", "details": err.Error()})
		return
	}

	// Save the updated Seeker record
	seeker.ProfessionalSummary = updatedProfessionalSummaryJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to patch professional summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Professional summary updated"})
}
