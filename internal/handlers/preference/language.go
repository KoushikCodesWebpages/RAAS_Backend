package preference

import (
	"RAAS/core/config"
	"RAAS/internal/dto"
	"RAAS/internal/handlers/repository"
	"RAAS/internal/models"


	"context"
	"log"
	"net/http"
	"strconv"
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
	mediaUploadHandler := repository.NewMediaUploadHandler(repository.GetBlobServiceClient())

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
	if err := repository.AppendToLanguages(&seeker, input, fileURL); err != nil {
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

	languages, err := repository.GetLanguages(&seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing languages"})
		log.Printf("Error processing languages for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
	})
}

// UpdateLanguage handles the update of a language entry with file upload
func (h *LanguageHandler) UpdateLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	db := c.MustGet("db").(*mongo.Database)
	seekersCollection := db.Collection("seekers")

	// Get the language index from the URL parameter
	id := c.Param("id")
	index, err := strconv.Atoi(id)
	if err != nil || index <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language index. Must be a positive integer."})
		return
	}

	// Bind the input data for the language
	var input dto.LanguageRequest
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		log.Printf("Error binding input: %v", err)
		return
	}

	// Handle file upload (if present)
	var fileURL string
	mediaUploadHandler := repository.NewMediaUploadHandler(repository.GetBlobServiceClient())

	// Check if a file was uploaded
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

	// Check if the index is valid in the current list of languages
	if index > len(seeker.Languages) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language index out of range"})
		return
	}

	// Prepare the updated language
	updatedLanguage := bson.M{
		"language":      input.LanguageName,
		"proficiency":   input.ProficiencyLevel,
		"certificate_file": fileURL, // Assuming fileURL is for the certificate
	}

	// Update the language entry at the specified index
	seeker.Languages[index-1] = updatedLanguage

	// Update the seeker document in the database
	update := bson.M{
		"$set": bson.M{
			"languages": seeker.Languages,
		},
	}

	updateResult, err := seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated language"})
		log.Printf("Failed to update language for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching seeker found to update"})
		log.Printf("No matching seeker found for auth_user_id: %s", userID)
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Language updated successfully",
	})
}

// DeleteLanguage handles deleting an existing language entry
func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	db := c.MustGet("db").(*mongo.Database)
	seekersCollection := db.Collection("seekers")

	id := c.Param("id")

	index, err := strconv.Atoi(id)
	if err != nil || index <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language index"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	if err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve seeker"})
		}
		return
	}

	if index > len(seeker.Languages) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language index out of range"})
		return
	}

	// Remove the language entry at index-1
	seeker.Languages = append(seeker.Languages[:index-1], seeker.Languages[index:]...)

	update := bson.M{
		"$set": bson.M{
			"languages": seeker.Languages,
		},
	}

	if _, err := seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete language entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language deleted successfully"})
}
