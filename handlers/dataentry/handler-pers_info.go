package dataentry

import (
	"RAAS/models"
	"RAAS/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"net/http"
)

// PersonalInfoHandler struct
type PersonalInfoHandler struct {
	DB *gorm.DB
}

// NewPersonalInfoHandler creates a new handler instance
func NewPersonalInfoHandler(db *gorm.DB) *PersonalInfoHandler {
	return &PersonalInfoHandler{DB: db}
}

// CreatePersonalInfo handles creation of personal info
func (h *PersonalInfoHandler) CreatePersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.PersonalInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
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

	if err := h.DB.Create(&personalInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create personal info", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, personalInfo)
}

// GetPersonalInfo retrieves personal info of the authenticated user
func (h *PersonalInfoHandler) GetPersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var personalInfo models.PersonalInfo
	if err := h.DB.First(&personalInfo, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal info not found"})
		return
	}

	c.JSON(http.StatusOK, personalInfo)
}

// UpdatePersonalInfo performs full update on personal info
func (h *PersonalInfoHandler) UpdatePersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.PersonalInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var personalInfo models.PersonalInfo
	if err := h.DB.First(&personalInfo, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal info not found"})
		return
	}

	personalInfo.FirstName = input.FirstName
	personalInfo.SecondName = input.SecondName
	personalInfo.DateOfBirth = input.DateOfBirth
	personalInfo.Address = input.Address
	personalInfo.LinkedInProfile = input.LinkedInProfile

	if err := h.DB.Save(&personalInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, personalInfo)
}

// PatchPersonalInfo handles partial updates to personal info
func (h *PersonalInfoHandler) PatchPersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Only allow specific fields
	allowedFields := map[string]string{
		"firstName":       "first_name",
		"secondName":      "second_name",
		"dob":             "date_of_birth",
		"address":         "address",
		"linkedinProfile": "linked_in_profile",
	}

	validUpdates := make(map[string]interface{})
	for k, v := range updates {
		if dbField, ok := allowedFields[k]; ok {
			validUpdates[dbField] = v
		}
	}

	if len(validUpdates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	if err := h.DB.Model(&models.PersonalInfo{}).
		Where("auth_user_id = ?", userID).
		Updates(validUpdates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to patch personal info", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Personal info updated"})
}
