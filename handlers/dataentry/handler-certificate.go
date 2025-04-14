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

type CertificateHandler struct {
	DB *gorm.DB
}

func NewCertificateHandler(db *gorm.DB) *CertificateHandler {
	return &CertificateHandler{DB: db}
}
func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)

    // Validate the file type
	mediaUploadHandler := features.NewMediaUploadHandler(features.GetBlobServiceClient(), config.Cfg.AzureBlobContainer)

    _, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file", "details": err.Error()})
        return
    }

    if !mediaUploadHandler.ValidateFileType(header) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
        return
    }

    // Upload the file to Azure Blob Storage
    fileURL, err := mediaUploadHandler.UploadMedia(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
        return
    }

    // Get the certificate details from the request
    certificateName := c.PostForm("CertificateName")
    certificateNumber := c.PostForm("CertificateNumber")

    if certificateName == "" || certificateNumber == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Certificate name and number are required"})
        return
    }

    certificate := models.Certificate{
        AuthUserID:      userID,
        CertificateName: certificateName,
        CertificateFile: fileURL, // Save the URL of the file
        CertificateNumber: certificateNumber,
    }

    // Create the certificate record in DB
    tx := h.DB.Begin()
    if err := tx.Create(&certificate).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create certificate", "details": err.Error()})
        return
    }

    // Update the UserEntryTimeline
    var timeline models.UserEntryTimeline
    if err := tx.First(&timeline, "user_id = ?", userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
        return
    }

    timeline.CertificatesCompleted = true
    if err := tx.Save(&timeline).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timeline", "details": err.Error()})
        return
    }

    tx.Commit()

    // Send back the certificate details in the response
    c.JSON(http.StatusCreated, dto.CertificateResponse{
        ID:                certificate.ID,
        AuthUserID:        certificate.AuthUserID,
        CertificateName:   certificate.CertificateName,
        CertificateFile:   certificate.CertificateFile,
        CertificateNumber: &certificate.CertificateNumber,
    })
}

func (h *CertificateHandler) GetCertificate(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)
    id := c.Param("id")

    var certificate models.Certificate
    if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).First(&certificate).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
        return
    }

    c.JSON(http.StatusOK, dto.CertificateResponse{
        ID:                certificate.ID,
        AuthUserID:        certificate.AuthUserID,
        CertificateName:   certificate.CertificateName,
        CertificateFile:   certificate.CertificateFile,
        CertificateNumber: &certificate.CertificateNumber,
    })
}


// PutCertificate updates an existing certificate record for the authenticated user
func (h *CertificateHandler) PutCertificate(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var existing models.Certificate
	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).First(&existing).Error; err != nil {
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

	if err := h.DB.Save(&updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update certificate", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate updated"})
}

// DeleteCertificate deletes an existing certificate record for the authenticated user
func (h *CertificateHandler) DeleteCertificate(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).Delete(&models.Certificate{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete certificate", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted"})
}