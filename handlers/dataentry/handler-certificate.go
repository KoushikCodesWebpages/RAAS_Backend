package dataentry

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"RAAS/handlers/features"
	"os"
	// "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	// "github.com/joho/godotenv"
	// "os"
	// "path/filepath"
	// "fmt"
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
    mediaUploadHandler := features.NewMediaUploadHandler(features.GetBlobServiceClient(), os.Getenv("AZURE_BLOB_CONTAINER"))
    _, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
        return
    }

    if !mediaUploadHandler.ValidateFileType(header) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
        return
    }

    // Upload the file to Azure Blob Storage
    fileURL, err := mediaUploadHandler.UploadMedia(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
        return
    }

    // Get the certificate details from the request
    certificateName := c.PostForm("CertificateName")
    certificateNumber := c.PostForm("CertificateNumber")

    certificate := models.Certificate{
        AuthUserID:      userID,
        CertificateName: certificateName,
        CertificateFile: fileURL, // Save the URL of the file
        CertificateNumber: certificateNumber,
    }

    // Create the certificate record in DB
    if err := h.DB.Create(&certificate).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create certificate", "details": err.Error()})
        return
    }

    // Update the UserEntryTimeline
    var timeline models.UserEntryTimeline
    if err := h.DB.First(&timeline, "user_id = ?", userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
        return
    }

    timeline.CertificatesCompleted = true
    if err := h.DB.Save(&timeline).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timeline", "details": err.Error()})
        return
    }

    // Send back the certificate details in the response
    c.JSON(http.StatusCreated, dto.CertificateResponse{
        ID:                certificate.ID,
        AuthUserID:        certificate.AuthUserID,
        CertificateName:   certificate.CertificateName,
        CertificateFile:   certificate.CertificateFile,
        CertificateNumber: &certificate.CertificateNumber,
    })
}
// GetCertificates retrieves all certificate records for the authenticated user
func (h *CertificateHandler) GetCertificates(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var certificates []models.Certificate
	if err := h.DB.Where("auth_user_id = ?", userID).Find(&certificates).Error; err != nil {
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