package dataentry

import (

	"log"
	"net/http"
    "encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

    "RAAS/config"
	"RAAS/dto"
	"RAAS/handlers/features"
	"RAAS/models"
)

type LanguageHandler struct {
	DB *gorm.DB
}

func NewLanguageHandler(db *gorm.DB) *LanguageHandler {
	return &LanguageHandler{DB: db}
}

// CreateLanguage adds a new language record for the authenticated user
func (h *LanguageHandler) CreateLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	// Retrieve language name and proficiency from form fields
	languageName := c.PostForm("LanguageName")
	proficiency := c.PostForm("ProficiencyLevel")

	if languageName == "" || proficiency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language name and proficiency are required"})
		return
	}

	// Handle file upload
	mediaUploadHandler := features.NewMediaUploadHandler(features.GetBlobServiceClient())

	_, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("[WARN] No file uploaded: %v", err)
	}

	var fileURL string
	if header != nil {
		if !mediaUploadHandler.ValidateFileType(header) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		fileURL, err = mediaUploadHandler.UploadMedia(c, config.Cfg.AzureLanguagesContainer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}
	}

	// Fetch seeker profile
	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker profile not found"})
		return
	}

	// Unmarshal existing languages
	var languages []map[string]interface{}
	if len(seeker.Languages) > 0 {
		if err := json.Unmarshal(seeker.Languages, &languages); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse languages", "details": err.Error()})
			return
		}
	}

	// Append new language
	newLanguage := map[string]interface{}{
		"language":     languageName,
		"proficiency":  proficiency,
		"certificateFile": fileURL,
	}
	languages = append(languages, newLanguage)

	// Marshal back to JSON
	updatedLanguagesJSON, err := json.Marshal(languages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated languages", "details": err.Error()})
		return
	}

	seeker.Languages = updatedLanguagesJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker languages", "details": err.Error()})
		return
	}

	// Return response
	response := dto.LanguageResponse{
		ID:               uint(len(languages)), // pseudo ID based on index
		AuthUserID:       userID,
		LanguageName:     languageName,
		CertificateFile:  fileURL,
		ProficiencyLevel: proficiency,
	}

	c.JSON(http.StatusCreated, response)
}
func (h *LanguageHandler) GetLanguages(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	if len(seeker.Languages) == 0 || string(seeker.Languages) == "null" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No language records found"})
		return
	}

	var languages []map[string]interface{}
	if err := json.Unmarshal(seeker.Languages, &languages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse language records"})
		return
	}

	var response []dto.LanguageResponse
	for idx, lang := range languages {
		response = append(response, dto.LanguageResponse{
			ID:              uint(idx + 1),
			AuthUserID:      userID,
			LanguageName:        lang["language"].(string),
			ProficiencyLevel:     lang["proficiency"].(string),
			CertificateFile: lang["certificateFile"].(string),
		})
	}

	c.JSON(http.StatusOK, response)
}


// UpdateLanguage updates an existing language record for the authenticated user
func (h *LanguageHandler) UpdateLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var existing models.Language
	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
		return
	}

	var updated models.Language
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Ensure these critical fields are preserved
	updated.ID = existing.ID
	updated.AuthUserID = userID
	updated.CreatedAt = existing.CreatedAt

	if err := h.DB.Save(&updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update language", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language updated"})
}

// DeleteLanguage deletes an existing language record for the authenticated user
func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).Delete(&models.Language{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete language", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language deleted"})
}
