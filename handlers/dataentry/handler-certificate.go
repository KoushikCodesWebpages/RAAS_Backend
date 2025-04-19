package dataentry

import (
	"RAAS/config"
	"RAAS/dto"
	"RAAS/handlers/features"
	"RAAS/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	certificateName := c.PostForm("CertificateName")
	certificateNumber := c.PostForm("CertificateNumber")

	if certificateName == "" || certificateNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Certificate name and number are required"})
		return
	}

	// File upload
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
		fileURL, err = mediaUploadHandler.UploadMedia(c, config.Cfg.AzureCertificatesContainer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}
	}

	// Get Seeker
	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	// Parse existing certificates
	var certificates []map[string]interface{}
	if len(seeker.Certificates) > 0 {
		if err := json.Unmarshal(seeker.Certificates, &certificates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse certificates", "details": err.Error()})
			return
		}
	}

	// Add new certificate
	newCert := map[string]interface{}{
		"certificateName":   certificateName,
		"certificateNumber": certificateNumber,
		"certificateFile":   fileURL,
	}
	certificates = append(certificates, newCert)

	updatedJSON, err := json.Marshal(certificates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal certificates"})
		return
	}
	seeker.Certificates = updatedJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save seeker certificates"})
		return
	}

	response := dto.CertificateResponse{
		ID:                uint(len(certificates)),
		AuthUserID:        userID,
		CertificateName:   certificateName,
		CertificateNumber: &certificateNumber,
		CertificateFile:   fileURL,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CertificateHandler) GetCertificates(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	if len(seeker.Certificates) == 0 || string(seeker.Certificates) == "null" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No certificate records found"})
		return
	}

	var certs []map[string]interface{}
	if err := json.Unmarshal(seeker.Certificates, &certs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse certificate records"})
		return
	}

	var response []dto.CertificateResponse
	for idx, cert := range certs {
		num := cert["certificateNumber"].(string)
		response = append(response, dto.CertificateResponse{
			ID:                uint(idx + 1),
			AuthUserID:        userID,
			CertificateName:   cert["certificateName"].(string),
			CertificateNumber: &num,
			CertificateFile:   cert["certificateFile"].(string),
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *CertificateHandler) PatchCertificate(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var updateFields map[string]interface{}
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	var certs []map[string]interface{}
	if err := json.Unmarshal(seeker.Certificates, &certs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse certificates"})
		return
	}

	index, err := strconv.Atoi(id)
	if err != nil || index <= 0 || index > len(certs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate index"})
		return
	}

	entry := certs[index-1]
	for key, value := range updateFields {
		if _, exists := entry[key]; exists {
			entry[key] = value
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid field: %s", key)})
			return
		}
	}
	certs[index-1] = entry

	updatedJSON, err := json.Marshal(certs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated certificates"})
		return
	}

	seeker.Certificates = updatedJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker certificates"})
		return
	}

	num := entry["certificateNumber"].(string)
	c.JSON(http.StatusOK, dto.CertificateResponse{
		ID:                uint(index),
		AuthUserID:        userID,
		CertificateName:   entry["certificateName"].(string),
		CertificateNumber: &num,
		CertificateFile:   entry["certificateFile"].(string),
	})
}

func (h *CertificateHandler) DeleteCertificate(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	var certs []map[string]interface{}
	if err := json.Unmarshal(seeker.Certificates, &certs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse certificates"})
		return
	}

	index, err := strconv.Atoi(id)
	if err != nil || index <= 0 || index > len(certs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate index"})
		return
	}

	certs = append(certs[:index-1], certs[index:]...)

	updatedJSON, err := json.Marshal(certs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated certificates"})
		return
	}

	seeker.Certificates = updatedJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted successfully"})
}
