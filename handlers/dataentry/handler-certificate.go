package dataentry

import (
	"RAAS/config"
	"RAAS/dto"
	"RAAS/handlers/features"
	"RAAS/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CertificateHandler struct {
	DB *gorm.DB
}

func NewCertificateHandler(db *gorm.DB) *CertificateHandler {
	return &CertificateHandler{DB: db}
}
func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)
    log.Printf("[DEBUG] Starting certificate creation for user: %s", userID)

    mediaUploadHandler := features.NewMediaUploadHandler(features.GetBlobServiceClient())
    log.Printf("[DEBUG] Initialized MediaUploadHandler")

    _, header, err := c.Request.FormFile("file")
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve file from request: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file", "details": err.Error()})
        return
    }

    log.Printf("[DEBUG] Received file: %s", header.Filename)

    if !mediaUploadHandler.ValidateFileType(header) {
        log.Printf("[ERROR] Invalid file type: %s", header.Filename)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
        return
    }

    fileURL, err := mediaUploadHandler.UploadMedia(c, config.Cfg.AzureCertificatesContainer)
    if err != nil {
        log.Printf("[ERROR] Failed to upload file to Azure Blob Storage: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
        return
    }

    log.Printf("[DEBUG] File uploaded successfully. URL: %s", fileURL)

    certificateName := c.PostForm("CertificateName")
    certificateNumber := c.PostForm("CertificateNumber")
    log.Printf("[DEBUG] Certificate form data - Name: %s, Number: %s", certificateName, certificateNumber)

    if certificateName == "" || certificateNumber == "" {
        log.Printf("[ERROR] Missing certificate name or number")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Certificate name and number are required"})
        return
    }

    certificate := models.Certificate{
        AuthUserID:        userID,
        CertificateName:   certificateName,
        CertificateFile:   fileURL,
        CertificateNumber: certificateNumber,
    }

    tx := h.DB.Begin()
    log.Printf("[DEBUG] Database transaction started")

    if err := tx.Create(&certificate).Error; err != nil {
        tx.Rollback()
        log.Printf("[ERROR] Failed to insert certificate: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create certificate", "details": err.Error()})
        return
    }

    log.Printf("[DEBUG] Certificate record created with ID: %s", certificate.CertificateNumber)

    var timeline models.UserEntryTimeline
    if err := tx.First(&timeline, "user_id = ?", userID).Error; err != nil {
        tx.Rollback()
        log.Printf("[ERROR] User entry timeline not found for user: %s", userID)
        c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
        return
    }

    timeline.CertificatesCompleted = true
    if err := tx.Save(&timeline).Error; err != nil {
        tx.Rollback()
        log.Printf("[ERROR] Failed to update user entry timeline: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timeline", "details": err.Error()})
        return
    }

    tx.Commit()
    log.Printf("[DEBUG] Transaction committed successfully")

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