package preference

import (
	"RAAS/core/config"
	"RAAS/internal/dto"
	"RAAS/internal/models"
	"RAAS/internal/handlers/repository"

	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateHandler struct{}

func NewCertificateHandler() *CertificateHandler {
	return &CertificateHandler{}
}

func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	db := c.MustGet("db").(*mongo.Database)
	seekersCollection := db.Collection("seekers")
	entryTimelineCollection := db.Collection("user_entry_timelines")

	var input dto.CertificateRequest
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		log.Printf("Error binding input: %v", err)
		return
	}

	// Optional file upload
	var fileURL string
	mediaUploadHandler := repository.NewMediaUploadHandler(repository.GetBlobServiceClient())

	_, header, err := c.Request.FormFile("file")
	if err == nil && header != nil {
		if !mediaUploadHandler.ValidateFileType(header) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}
		fileURL, err = mediaUploadHandler.UploadMedia(c, config.Cfg.Cloud.AzureCertificatesContainer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}
	} else {
		// File is optional — just log and continue
		fileURL = ""
		log.Printf("[INFO] No certificate file uploaded for user %s. Proceeding without it.", userID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	if err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
			log.Printf("Seeker not found for auth_user_id: %s", userID)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving seeker"})
			log.Printf("Error retrieving seeker for auth_user_id: %s, Error: %v", userID, err)
		}
		return
	}

	// Append the new certificate
	if err := repository.AppendToCertificates(&seeker, input, fileURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process certificate"})
		log.Printf("Failed to process certificate for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	// Update seeker document
	update := bson.M{
		"$set": bson.M{
			"certificates": seeker.Certificates,
		},
	}

	updateResult, err := seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save certificate"})
		log.Printf("Failed to update certificate for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching seeker found to update"})
		log.Printf("No matching seeker found for auth_user_id: %s", userID)
		return
	}

	// Update user entry timeline to mark certificates completed
	timelineUpdate := bson.M{
		"$set": bson.M{
			"certificates_completed": true,
		},
	}
	if _, err := entryTimelineCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, timelineUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user entry timeline"})
		log.Printf("Failed to update user entry timeline for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Certificate added successfully",
	})
}

	
// GetCertificates handles the retrieval of a user's certificates
func (h *CertificateHandler) GetCertificates(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	db := c.MustGet("db").(*mongo.Database)
	seekersCollection := db.Collection("seekers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	if err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
			log.Printf("Seeker not found for auth_user_id: %s", userID)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving seeker"})
			log.Printf("Error retrieving seeker for auth_user_id: %s, Error: %v", userID, err)
		}
		return
	}

	if len(seeker.Certificates) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No certificates found"})
		return
	}

	certificates, err := repository.GetCertificates(&seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing certificates"})
		log.Printf("Error processing certificates for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"certificates": certificates,
	})
}

// func (h *CertificateHandler) PatchCertificate(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)
// 	id := c.Param("id")

// 	var updateFields map[string]interface{}
// 	if err := c.ShouldBindJSON(&updateFields); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
// 		return
// 	}

// 	var seeker models.Seeker
// 	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	var certs []map[string]interface{}
// 	if err := json.Unmarshal(seeker.Certificates, &certs); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse certificates"})
// 		return
// 	}

// 	index, err := strconv.Atoi(id)
// 	if err != nil || index <= 0 || index > len(certs) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate index"})
// 		return
// 	}

// 	entry := certs[index-1]
// 	for key, value := range updateFields {
// 		if _, exists := entry[key]; exists {
// 			entry[key] = value
// 		} else {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid field: %s", key)})
// 			return
// 		}
// 	}
// 	certs[index-1] = entry

// 	updatedJSON, err := json.Marshal(certs)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated certificates"})
// 		return
// 	}

// 	seeker.Certificates = updatedJSON
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker certificates"})
// 		return
// 	}

// 	num := entry["certificateNumber"].(string)
// 	c.JSON(http.StatusOK, dto.CertificateResponse{
// 		ID:                uint(index),
// 		AuthUserID:        userID,
// 		CertificateName:   entry["certificateName"].(string),
// 		CertificateNumber: &num,
// 		CertificateFile:   entry["certificateFile"].(string),
// 	})
// }

// func (h *CertificateHandler) DeleteCertificate(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)
// 	id := c.Param("id")

// 	var seeker models.Seeker
// 	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	var certs []map[string]interface{}
// 	if err := json.Unmarshal(seeker.Certificates, &certs); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse certificates"})
// 		return
// 	}

// 	index, err := strconv.Atoi(id)
// 	if err != nil || index <= 0 || index > len(certs) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate index"})
// 		return
// 	}

// 	certs = append(certs[:index-1], certs[index:]...)

// 	updatedJSON, err := json.Marshal(certs)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated certificates"})
// 		return
// 	}

// 	seeker.Certificates = updatedJSON
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted successfully"})
// }
