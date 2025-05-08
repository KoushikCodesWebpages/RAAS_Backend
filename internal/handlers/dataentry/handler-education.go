package dataentry

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"RAAS/internal/dto"
	"RAAS/internal/models"
	"RAAS/internal/handlers"
)

type EducationHandler struct{}

func NewEducationHandler() *EducationHandler {
	return &EducationHandler{}
}

// CreateEducation handles the creation or update of a single education entry
func (h *EducationHandler) CreateEducation(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	db := c.MustGet("db").(*mongo.Database)
	seekersCollection := db.Collection("seekers")
	entryTimelineCollection := db.Collection("user_entry_timelines")

	var input dto.EducationRequest
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

	// Create an EducationRequest from the input
	education := dto.EducationRequest{
		Degree:       input.Degree,
		Institution:  input.Institution,
		FieldOfStudy: input.FieldOfStudy,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		Achievements: input.Achievements,
	}

	// Use AppendToEducation to add the new education
	if err := handlers.AppendToEducation(&seeker, education); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process education"})
		log.Printf("Failed to process education for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	// Update seeker document with the new education
	update := bson.M{
		"$set": bson.M{
			"education": seeker.Education, // Save updated education records
		},
	}

	updateResult, err := seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save education"})
		log.Printf("Failed to update education for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching seeker found to update"})
		log.Printf("No matching seeker found for auth_user_id: %s", userID)
		return
	}

	// Update user entry timeline to mark education completed
	timelineUpdate := bson.M{
		"$set": bson.M{
			"educations_completed": true,
		},
	}

	if _, err := entryTimelineCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, timelineUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user entry timeline"})
		log.Printf("Failed to update user entry timeline for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Education added successfully",
	})
}
// GetEducationHandler handles the retrieval of a user's education records
func (h *EducationHandler) GetEducation(c *gin.Context) {
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

    // Check if the user has any education records
    if len(seeker.Education) == 0 {
        c.JSON(http.StatusNoContent, gin.H{"message": "No education records found"})
        return
    }

    // Fetch the education data (could be a function similar to GetWorkExperience)
    educations, err := handlers.GetEducation(&seeker)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing education records"})
        log.Printf("Error processing education records for auth_user_id: %s, Error: %v", userID, err)
        return
    }

    // Return the education data in the response
    c.JSON(http.StatusOK, gin.H{
        "educations": educations,
    })
}


// func (h *EducationHandler) PatchEducation(c *gin.Context) {
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

// 	var educations []map[string]interface{}
// 	if err := json.Unmarshal(seeker.Educations, &educations); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse educations"})
// 		return
// 	}

// 	index, err := strconv.Atoi(id)
// 	if err != nil || index <= 0 || index > len(educations) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid education index"})
// 		return
// 	}

// 	// Apply updates
// 	entry := educations[index-1]
// 	for key, value := range updateFields {
// 		if _, exists := entry[key]; exists {
// 			entry[key] = value
// 		} else {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid field: %s", key)})
// 			return
// 		}
// 	}
// 	educations[index-1] = entry

// 	updatedJSON, err := json.Marshal(educations)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated educations"})
// 		return
// 	}

// 	seeker.Educations = updatedJSON
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
// 	response := dto.EducationResponse{
// 		ID:             uint(index),
// 		AuthUserID:     userID,
// 		Degree:         entry["degree"].(string),
// 		Institution:    entry["institution"].(string),
// 		FieldOfStudy:   entry["fieldOfStudy"].(string),
// 		StartDate:      startDate,
// 		EndDate:        endDate,
// 		Achievements:   entry["achievements"].(string),
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// func (h *EducationHandler) DeleteEducation(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)
// 	id := c.Param("id")

// 	var seeker models.Seeker
// 	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	var educations []map[string]interface{}
// 	if err := json.Unmarshal(seeker.Educations, &educations); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse educations"})
// 		return
// 	}

// 	index, err := strconv.Atoi(id)
// 	if err != nil || index <= 0 || index > len(educations) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid education index"})
// 		return
// 	}

// 	// Remove the education at the specified index (index - 1 since it's 1-based in API)
// 	educations = append(educations[:index-1], educations[index:]...)

// 	updatedJSON, err := json.Marshal(educations)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated educations"})
// 		return
// 	}

// 	seeker.Educations = updatedJSON
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Education deleted successfully"})
// }
