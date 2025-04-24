package dataentry

import (

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"

	"RAAS/models"
	"RAAS/dto"
	"RAAS/handlers"
)

// PersonalInfoHandler struct
type PersonalInfoHandler struct {
	DB *gorm.DB
}

// NewPersonalInfoHandler creates a new handler instance
func NewPersonalInfoHandler(db *gorm.DB) *PersonalInfoHandler {
	return &PersonalInfoHandler{DB: db}
}

// CreatePersonalInfo handles creation of personal info within Seeker model
func (h *PersonalInfoHandler) CreatePersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.PersonalInfoRequest
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

	// Check if PersonalInfo is already filled
	if len(seeker.PersonalInfo) > 0 {
		var existingInfo map[string]interface{}
		if err := json.Unmarshal(seeker.PersonalInfo, &existingInfo); err == nil && len(existingInfo) > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Personal info already filled"})
			return
		}
	}

	if input.DateOfBirth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date of birth cannot be empty"})
		return
	}

	// Use SetPersonalInfo utility function
	if err := handlers.SetPersonalInfo(&seeker, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set personal info", "details": err.Error()})
		return
	}

	// Save the updated Seeker record
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
		return
	}

	// Return the PersonalInfoResponse
	personalInfoResponse := dto.PersonalInfoResponse{
		AuthUserID:      seeker.AuthUserID,
		FirstName:       input.FirstName,
		SecondName:      input.SecondName,
		DateOfBirth:     input.DateOfBirth,
		Address:         input.Address,
		LinkedInProfile: input.LinkedInProfile,
	}

	c.JSON(http.StatusCreated, personalInfoResponse)
}

// GetPersonalInfo retrieves personal info of the authenticated user
func (h *PersonalInfoHandler) GetPersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	if len(seeker.PersonalInfo) == 0 || string(seeker.PersonalInfo) == "null" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal info not yet filled"})
		return
	}

	// Use GetPersonalInfo utility function
	personalInfo, err := handlers.GetPersonalInfo(&seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse personal info", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, personalInfo)
}

// UpdatePersonalInfo handles the update of personal info for Seeker model
func (h *PersonalInfoHandler) UpdatePersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.PersonalInfoRequest
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

	// Use SetPersonalInfo utility function to update
	if err := handlers.SetPersonalInfo(&seeker, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
		return
	}

	// Save the updated Seeker record
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated personal info", "details": err.Error()})
		return
	}

	// Return the updated PersonalInfoResponse
	personalInfoResponse := dto.PersonalInfoResponse{
		AuthUserID:      seeker.AuthUserID,
		FirstName:       input.FirstName,
		SecondName:      input.SecondName,
		DateOfBirth:     input.DateOfBirth,
		Address:         input.Address,
		LinkedInProfile: input.LinkedInProfile,
	}

	c.JSON(http.StatusOK, personalInfoResponse)
}

// PatchPersonalInfo handles patching personal info for Seeker model
func (h *PersonalInfoHandler) PatchPersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.PersonalInfoRequest
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

	// Get the current personal info
	personalInfo, err := handlers.GetPersonalInfo(&seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse personal info", "details": err.Error()})
		return
	}

	// Check if DateOfBirth or Address are being updated and return an error if so
	if input.DateOfBirth != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Updating DateOfBirth is not allowed"})
		return
	}
	if input.Address != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Updating Address is not allowed"})
		return
	}

	// Update only the allowed fields
	if input.FirstName != "" {
		personalInfo.FirstName = input.FirstName
	}
	if input.SecondName != nil {
		personalInfo.SecondName = input.SecondName
	}
	if input.LinkedInProfile != nil {
		personalInfo.LinkedInProfile = input.LinkedInProfile
	}

	// Use SetPersonalInfo utility function to update the personal info
	if err := handlers.SetPersonalInfo(&seeker, personalInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
		return
	}

	// Save the updated Seeker record
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated personal info", "details": err.Error()})
		return
	}

	// Return the updated PersonalInfoResponse
	personalInfoResponse := dto.PersonalInfoResponse{
		AuthUserID:      seeker.AuthUserID,
		FirstName:       personalInfo.FirstName,
		SecondName:      personalInfo.SecondName,  // Nullable, so it could be null
		DateOfBirth:     personalInfo.DateOfBirth, // Unchanged
		Address:         personalInfo.Address,     // Unchanged
		LinkedInProfile: personalInfo.LinkedInProfile,  // Nullable, so it could be null
	}

	c.JSON(http.StatusOK, personalInfoResponse)
}
