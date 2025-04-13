package dataentry

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LanguageHandler struct
type LanguageHandler struct {
	DB *gorm.DB
}

// NewLanguageHandler creates a new LanguageHandler
func NewLanguageHandler(db *gorm.DB) *LanguageHandler {
	return &LanguageHandler{DB: db}
}

// CreateLanguage creates a new language entry for the authenticated user
func (h *LanguageHandler) CreateLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.LanguageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	lang := models.Language{
		AuthUserID:       userID,
		LanguageName:     input.LanguageName,
		ProficiencyLevel: input.ProficiencyLevel,
	}

	if input.CertificateFile != nil {
		lang.CertificateFile = *input.CertificateFile
	}

	if err := h.DB.Create(&lang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create language entry", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.LanguageResponse{
		ID:               lang.ID,
		AuthUserID:       userID,
		LanguageName:     lang.LanguageName,
		CertificateFile:  input.CertificateFile,
		ProficiencyLevel: lang.ProficiencyLevel,
	})
}

// GetLanguages retrieves all language entries for the authenticated user
func (h *LanguageHandler) GetLanguages(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var langs []models.Language
	if err := h.DB.Where("auth_user_id = ?", userID).Find(&langs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch language entries", "details": err.Error()})
		return
	}

	var response []dto.LanguageResponse
	for _, lang := range langs {
		file := lang.CertificateFile
		response = append(response, dto.LanguageResponse{
			ID:               lang.ID,
			AuthUserID:       lang.AuthUserID,
			LanguageName:     lang.LanguageName,
			CertificateFile:  &file,
			ProficiencyLevel: lang.ProficiencyLevel,
		})
	}

	c.JSON(http.StatusOK, response)
}

// PutLanguage updates an existing language entry for the authenticated user
func (h *LanguageHandler) PutLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var existing models.Language
	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Language entry not found"})
		return
	}

	var updated models.Language
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Preserve non-editable fields
	updated.ID = existing.ID
	updated.AuthUserID = userID
	updated.CreatedAt = existing.CreatedAt

	if err := h.DB.Save(&updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update language", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language updated successfully"})
}

// DeleteLanguage deletes an existing language entry for the authenticated user
func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).Delete(&models.Language{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete language", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language deleted"})
}
