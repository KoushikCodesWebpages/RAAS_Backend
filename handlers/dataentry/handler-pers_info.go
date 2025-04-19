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
	
	// Prepare the personal info map
	personalInfo := map[string]interface{}{
		"firstName":       input.FirstName,
		"secondName":      input.SecondName,
		"dateOfBirth":     input.DateOfBirth,
		"address":         input.Address,
		"linkedInProfile": input.LinkedInProfile,
	}

	// Marshal the personal info into JSON
	personalInfoJSON, err := json.Marshal(personalInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal personal info", "details": err.Error()})
		return
	}

	// Update Seeker's PersonalInfo field
	seeker.PersonalInfo = personalInfoJSON

	// Save the updated Seeker record
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
		return
	}

	// Update the UserEntryTimeline - set PersonalInfosCompleted to true
	var timeline models.UserEntryTimeline
	if err := h.DB.First(&timeline, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
		return
	}

	timeline.PersonalInfosCompleted = true

	// Save the updated timeline
	if err := h.DB.Save(&timeline).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timeline", "details": err.Error()})
		return
	}

	// Return only the PersonalInfoResponse
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
// GetPersonalInfo retrieves the personal information of the authenticated user
func (h *PersonalInfoHandler) GetPersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Check if PersonalInfo is empty or uninitialized
	if len(seeker.PersonalInfo) == 0 || string(seeker.PersonalInfo) == "null" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal info not yet filled"})
		return
	}

	// Unmarshal the PersonalInfo JSON field into a map
	var personalInfo map[string]interface{}
	if err := json.Unmarshal(seeker.PersonalInfo, &personalInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse personal info", "details": err.Error()})
		return
	}

	// Ensure all required fields exist and are of the correct type
	firstName, firstNameOk := personalInfo["first_name"].(string)
	secondName, secondNameOk := personalInfo["second_name"].(*string)
	dateOfBirth, dateOfBirthOk := personalInfo["date_of_birth"].(string)
	address, addressOk := personalInfo["address"].(string)
	linkedinProfile, linkedinProfileOk := personalInfo["linkedin_profile"].(*string)

	if !firstNameOk || !secondNameOk || !dateOfBirthOk || !addressOk || !linkedinProfileOk {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid personal info data",
			"details": "Expected 'first_name', 'second_name', 'date_of_birth', 'address', and 'linkedin_profile' to be strings",
		})
		return
	}

	// Create a response containing the PersonalInfo data
	personalInfoResponse := dto.PersonalInfoResponse{
		AuthUserID:      seeker.AuthUserID,
		FirstName:       firstName,
		SecondName:      secondName,
		DateOfBirth:     dateOfBirth,
		Address:         address,
		LinkedInProfile: linkedinProfile,
	}

	c.JSON(http.StatusOK, personalInfoResponse)
}

// UpdatePersonalInfo performs full update on personal info within Seeker model
func (h *PersonalInfoHandler) UpdatePersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.PersonalInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Prepare the personal info map
	personalInfo := map[string]interface{}{
		"firstName":       input.FirstName,
		"secondName":      input.SecondName,
		"dateOfBirth":     input.DateOfBirth,
		"address":         input.Address,
		"linkedInProfile": input.LinkedInProfile,
	}

	// Marshal the personal info into JSON
	personalInfoJSON, err := json.Marshal(personalInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal personal info", "details": err.Error()})
		return
	}

	// Update Seeker's PersonalInfo field
	seeker.PersonalInfo = personalInfoJSON

	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seeker)
}

// PatchPersonalInfo handles partial updates to personal info within Seeker model
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

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Convert the PersonalInfo field back to a map (deserialize)
	var personalInfo map[string]interface{}
	if err := json.Unmarshal(seeker.PersonalInfo, &personalInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal personal info", "details": err.Error()})
		return
	}

	// Update the map with the valid updates	
	for key, value := range validUpdates {
		personalInfo[key] = value
	}

	// Marshal the updated personal info back to JSON
	updatedPersonalInfoJSON, err := json.Marshal(personalInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated personal info", "details": err.Error()})
		return
	}

	// Save the updated Seeker record
	seeker.PersonalInfo = updatedPersonalInfoJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to patch personal info", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Personal info updated"})
}
