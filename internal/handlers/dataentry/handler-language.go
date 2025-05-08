package dataentry

import (
	"RAAS/core/config"
	"RAAS/internal/dto"
	"RAAS/internal/handlers"
	"RAAS/internal/handlers/features"
	"RAAS/internal/models"


	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LanguageHandler struct{}

func NewLanguageHandler() *LanguageHandler {
	return &LanguageHandler{}
}

// CreateLanguage handles the creation or update of a single language entry
func (h *LanguageHandler) CreateLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	db := c.MustGet("db").(*mongo.Database)
	seekersCollection := db.Collection("seekers")
	entryTimelineCollection := db.Collection("user_entry_timelines")

	var input dto.LanguageRequest
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		log.Printf("Error binding input: %v", err)
		return
	}

	// Upload file if present
	var fileURL string
	mediaUploadHandler := features.NewMediaUploadHandler(features.GetBlobServiceClient())

	_, header, err := c.Request.FormFile("file")
	if err == nil && header != nil {
		if !mediaUploadHandler.ValidateFileType(header) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}
		fileURL, err = mediaUploadHandler.UploadMedia(c, config.Cfg.Cloud.AzureLanguagesContainer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}
	} else {
		log.Printf("[WARN] No file uploaded for language: %v", err)
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

	// Append the new language
	if err := handlers.AppendToLanguages(&seeker, input, fileURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process language"})
		log.Printf("Failed to process language for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	// Update seeker document
	update := bson.M{
		"$set": bson.M{
			"languages": seeker.Languages,
		},
	}

	updateResult, err := seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save language"})
		log.Printf("Failed to update language for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching seeker found to update"})
		log.Printf("No matching seeker found for auth_user_id: %s", userID)
		return
	}

	// Update user entry timeline to mark languages completed
	timelineUpdate := bson.M{
		"$set": bson.M{
			"languages_completed": true,
		},
	}
	if _, err := entryTimelineCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, timelineUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user entry timeline"})
		log.Printf("Failed to update user entry timeline for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Language added successfully",
	})
}

// GetLanguages handles the retrieval of a user's languages
func (h *LanguageHandler) GetLanguages(c *gin.Context) {
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

	if len(seeker.Languages) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No languages found"})
		return
	}

	languages, err := handlers.GetLanguages(&seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing languages"})
		log.Printf("Error processing languages for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
	})
}


// func (h *LanguageHandler) PatchLanguage(c *gin.Context) {
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

// 	var languages []map[string]interface{}
// 	if err := json.Unmarshal(seeker.Languages, &languages); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse languages"})
// 		return
// 	}

// 	index, err := strconv.Atoi(id)
// 	if err != nil || index <= 0 || index > len(languages) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language index"})
// 		return
// 	}

// 	entry := languages[index-1]
// 	for key, value := range updateFields {
// 		if _, exists := entry[key]; exists {
// 			entry[key] = value
// 		} else {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid field: %s", key)})
// 			return
// 		}
// 	}
// 	languages[index-1] = entry

// 	updatedJSON, err := json.Marshal(languages)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated languages"})
// 		return
// 	}

// 	seeker.Languages = updatedJSON
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
// 		return
// 	}

// 	response := dto.LanguageResponse{
// 		ID:              uint(index),
// 		AuthUserID:      userID,
// 		LanguageName:        entry["language"].(string),
// 		ProficiencyLevel:     entry["proficiency"].(string),
// 		CertificateFile: entry["certificateFile"].(string),
// 	}

// 	c.JSON(http.StatusOK, response)
// }


// func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)
// 	id := c.Param("id")

// 	var seeker models.Seeker
// 	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	var languages []map[string]interface{}
// 	if err := json.Unmarshal(seeker.Languages, &languages); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse languages"})
// 		return
// 	}

// 	index, err := strconv.Atoi(id)
// 	if err != nil || index <= 0 || index > len(languages) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language index"})
// 		return
// 	}

// 	// Remove the language at the specified index
// 	languages = append(languages[:index-1], languages[index:]...)

// 	updatedJSON, err := json.Marshal(languages)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated languages"})
// 		return
// 	}

// 	seeker.Languages = updatedJSON
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Language deleted successfully"})
// }
