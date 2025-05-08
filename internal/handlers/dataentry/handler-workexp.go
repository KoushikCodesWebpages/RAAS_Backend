package dataentry

import (

	"RAAS/internal/dto"
	"RAAS/internal/models"
	"RAAS/internal/handlers"

	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkExperienceHandler struct{}

func NewWorkExperienceHandler() *WorkExperienceHandler {
	return &WorkExperienceHandler{}
}

// CreateWorkExperience handles the creation or update of a single work experience
func (h *WorkExperienceHandler) CreateWorkExperience(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	db := c.MustGet("db").(*mongo.Database)
	seekersCollection := db.Collection("seekers")
	entryTimelineCollection := db.Collection("user_entry_timelines")

	var input dto.WorkExperienceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		log.Printf("Error binding input: %v", err)
		return
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

	// Create a dto.WorkExperienceRequest from the input
	workExperience := dto.WorkExperienceRequest{
		JobTitle:           input.JobTitle,
		CompanyName:        input.CompanyName,
		EmploymentType:     input.EmploymentType,
		StartDate:          input.StartDate,
		EndDate:            input.EndDate,
		KeyResponsibilities: input.KeyResponsibilities,
	}

	// Use AppendToWorkExperience to add the new experience
	if err := handlers.AppendToWorkExperience(&seeker, workExperience); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process work experience"})
		log.Printf("Failed to process work experience for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	// Update seeker document with the new work experiences
	update := bson.M{
		"$set": bson.M{
			"work_experiences": seeker.WorkExperiences, // Save updated work experiences
		},
	}

	updateResult, err := seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save work experience"})
		log.Printf("Failed to update work experience for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching seeker found to update"})
		log.Printf("No matching seeker found for auth_user_id: %s", userID)
		return
	}

	// Update user entry timeline to mark work experiences completed
	timelineUpdate := bson.M{
		"$set": bson.M{
			"work_experiences_completed": true,
		},
	}

	if _, err := entryTimelineCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, timelineUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user entry timeline"})
		log.Printf("Failed to update user entry timeline for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Work experience added successfully",
	})
}

// GetWorkExperienceHandler handles the retrieval of a user's work experiences
func (h *WorkExperienceHandler) GetWorkExperience(c *gin.Context) {
    // Extract user ID from the context
    userID := c.MustGet("userID").(string)
    db := c.MustGet("db").(*mongo.Database)
    seekersCollection := db.Collection("seekers")

    // Set up a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Find the seeker by their auth_user_id
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

    // Check if the user has work experiences
    if len(seeker.WorkExperiences) == 0 {
        c.JSON(http.StatusNoContent, gin.H{"message": "No work experiences found"})
        return
    }

    // Convert bson.M to the expected dto.WorkExperienceRequest type
    workExperiences, err := handlers.GetWorkExperience(&seeker)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing work experiences"})
        log.Printf("Error processing work experiences for auth_user_id: %s, Error: %v", userID, err)
        return
    }

    // Return the work experiences in the response
    c.JSON(http.StatusOK, gin.H{
        "work_experiences": workExperiences,
    })
}



// func (h *WorkExperienceHandler) PatchWorkExperience(c *gin.Context) {
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

// 	var workExperiences []map[string]interface{}
// 	if err := json.Unmarshal(seeker.WorkExperiences, &workExperiences); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse work experiences"})
// 		return
// 	}

// 	index, err := strconv.Atoi(id)
// 	if err != nil || index <= 0 || index > len(workExperiences) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work experience index"})
// 		return
// 	}

// 	// Apply updates
// 	entry := workExperiences[index-1]
// 	for key, value := range updateFields {
// 		if _, exists := entry[key]; exists {
// 			entry[key] = value
// 		} else {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid field: %s", key)})
// 			return
// 		}
// 	}
// 	workExperiences[index-1] = entry

// 	updatedJSON, err := json.Marshal(workExperiences)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated work experiences"})
// 		return
// 	}

// 	seeker.WorkExperiences = updatedJSON
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
// 		return
// 	}

// 	// Parse StartDate and EndDate (ensure both are valid)
// 	startDate, err := time.Parse("2006-01-02", entry["startDate"].(string))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid start date format"})
// 		return
// 	}

// 	// EndDate is required, so parse it directly without nil checks
// 	endDate, err := time.Parse("2006-01-02", entry["endDate"].(string))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid end date format"})
// 		return
// 	}

// 	// Create and return response with updated data
// 	response := dto.WorkExperienceResponse{
// 		ID:                  uint(index),
// 		AuthUserID:          userID,
// 		JobTitle:            entry["jobTitle"].(string),
// 		CompanyName:         entry["companyName"].(string),
// 		EmploymentType:      entry["employmentType"].(string),
// 		StartDate:           startDate,
// 		EndDate:             endDate,
// 		KeyResponsibilities: entry["keyResponsibilities"].(string),
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// func (h *WorkExperienceHandler) DeleteWorkExperience(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)
// 	id := c.Param("id")

// 	var seeker models.Seeker
// 	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	var workExperiences []map[string]interface{}
// 	if err := json.Unmarshal(seeker.WorkExperiences, &workExperiences); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse work experiences"})
// 		return
// 	}

// 	index, err := strconv.Atoi(id)
// 	if err != nil || index <= 0 || index > len(workExperiences) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work experience index"})
// 		return
// 	}

// 	// Remove the work experience at the specified index (index - 1 since it's 1-based in API)
// 	workExperiences = append(workExperiences[:index-1], workExperiences[index:]...)

// 	updatedJSON, err := json.Marshal(workExperiences)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated work experiences"})
// 		return
// 	}

// 	seeker.WorkExperiences = updatedJSON
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Work experience deleted successfully"})
// }



