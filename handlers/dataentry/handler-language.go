package dataentry

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"RAAS/handlers/features"
	"RAAS/config"
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

    languageName := c.PostForm("LanguageName")
    proficiencyLevel := c.PostForm("ProficiencyLevel")

    if languageName == "" || proficiencyLevel == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Language name and proficiency level are required"})
        return
    }

    var fileURL string
    _, header, err := c.Request.FormFile("file")
    if err == nil {
        // Initialize your media upload handler
        mediaUploadHandler := features.NewMediaUploadHandler(features.GetBlobServiceClient())

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

    language := models.Language{
        AuthUserID:       userID,
        LanguageName:     languageName,
        ProficiencyLevel: proficiencyLevel,
        CertificateFile:  fileURL,
    }

    tx := h.DB.Begin()

    if err := tx.Create(&language).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create language entry", "details": err.Error()})
        return
    }

    var timeline models.UserEntryTimeline
    if err := tx.First(&timeline, "user_id = ?", userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
        return
    }

    timeline.LanguagesCompleted = true
    if err := tx.Save(&timeline).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timeline", "details": err.Error()})
        return
    }

    tx.Commit()

    c.JSON(http.StatusCreated, dto.LanguageResponse{
        ID:               language.ID,
        AuthUserID:       language.AuthUserID,
        LanguageName:     language.LanguageName,
        CertificateFile:  &language.CertificateFile,
        ProficiencyLevel: language.ProficiencyLevel,
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
