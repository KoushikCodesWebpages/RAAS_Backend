package handlers

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateLanguage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
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

	if err := db.Create(&lang).Error; err != nil {
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

func GetLanguages(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var langs []models.Language
	if err := db.Where("auth_user_id = ?", userID).Find(&langs).Error; err != nil {
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

func PutLanguage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var existing models.Language
	if err := db.Where("id = ? AND auth_user_id = ?", id, userID).First(&existing).Error; err != nil {
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

	if err := db.Save(&updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update language", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language updated successfully"})
}


func DeleteLanguage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	if err := db.Where("id = ? AND auth_user_id = ?", id, userID).Delete(&models.Language{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete language", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language deleted"})
}
