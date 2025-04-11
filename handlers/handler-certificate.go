package handlers

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCertificate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.CertificateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	certificate := models.Certificate{
		AuthUserID:        userID,
		CertificateName:   input.CertificateName,
		CertificateFile:   input.CertificateFile,
	}

	if input.CertificateNumber != nil {
		certificate.CertificateNumber = *input.CertificateNumber
	}

	if err := db.Create(&certificate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create certificate", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CertificateResponse{
		ID:                certificate.ID,
		AuthUserID:        certificate.AuthUserID,
		CertificateName:   certificate.CertificateName,
		CertificateFile:   certificate.CertificateFile,
		CertificateNumber: input.CertificateNumber,
	})
}

func GetCertificates(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var certificates []models.Certificate
	if err := db.Where("auth_user_id = ?", userID).Find(&certificates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch certificates", "details": err.Error()})
		return
	}

	var response []dto.CertificateResponse
	for _, cert := range certificates {
		number := cert.CertificateNumber
		response = append(response, dto.CertificateResponse{
			ID:                cert.ID,
			AuthUserID:        cert.AuthUserID,
			CertificateName:   cert.CertificateName,
			CertificateFile:   cert.CertificateFile,
			CertificateNumber: &number,
		})
	}

	c.JSON(http.StatusOK, response)
}

func PutCertificate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var existing models.Certificate
	if err := db.Where("id = ? AND auth_user_id = ?", id, userID).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	var updated models.Certificate
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Ensure these critical fields are preserved
	updated.ID = existing.ID
	updated.AuthUserID = userID
	updated.CreatedAt = existing.CreatedAt

	if err := db.Save(&updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update certificate", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate updated"})
}


func DeleteCertificate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	if err := db.Where("id = ? AND auth_user_id = ?", id, userID).Delete(&models.Certificate{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete certificate", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted"})
}
