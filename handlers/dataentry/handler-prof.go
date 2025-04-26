package dataentry

import (
	"context"
	"log"
	"net/http"
	"time"

	"RAAS/dto"
	"RAAS/handlers"
	"RAAS/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfessionalSummaryHandler struct{}

func NewProfessionalSummaryHandler() *ProfessionalSummaryHandler {
	return &ProfessionalSummaryHandler{}
}

func (h *ProfessionalSummaryHandler) CreateProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	seekersCollection := c.MustGet("db").(*mongo.Database).Collection("seekers")
	entryTimelineCollection := c.MustGet("db").(*mongo.Database).Collection("user_entry_timelines")

	var input dto.ProfessionalSummaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		log.Printf("Error binding input: %s", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	// Fetch seeker by auth_user_id
	err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
			log.Printf("Seeker not found for auth_user_id: %s", userID)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving seeker"})
			log.Printf("Error retrieving seeker for auth_user_id: %s, Error: %s", userID, err.Error())
		}
		return
	}

	log.Printf("Found seeker: %+v", seeker)

	// Process and set professional summary
	if err := handlers.SetProfessionalSummary(&seeker, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process professional summary"})
		log.Printf("Failed to process professional summary for auth_user_id: %s, Error: %s", userID, err.Error())
		return
	}

	// Update professional summary in MongoDB (seekers collection)
	updateResult, err := seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, bson.M{"$set": bson.M{"professional_summary": seeker.ProfessionalSummary}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save professional summary"})
		log.Printf("Failed to update professional summary for auth_user_id: %s, Error: %s", userID, err.Error())
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching seeker found to update"})
		log.Printf("No matching seeker found for auth_user_id: %s", userID)
		return
	}

	message := "Professional summary created"
	if handlers.IsFieldFilled(seeker.ProfessionalSummary) {
		message = "Professional summary updated"
	}

	// Now, update the user entry progress in user_entry_timelines collection
	// Set professional_summaries_completed to true
	timelineUpdate := bson.M{
		"$set": bson.M{
			"professional_summaries_completed": true,
		},
	}

	timelineUpdateResult, err := entryTimelineCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, timelineUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user entry timeline"})
		log.Printf("Failed to update user entry timeline for auth_user_id: %s, Error: %s", userID, err.Error())
		return
	}

	log.Printf("Timeline update result: %+v", timelineUpdateResult)

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (h *ProfessionalSummaryHandler) GetProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	seekersCollection := c.MustGet("db").(*mongo.Database).Collection("seekers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
			log.Printf("Seeker not found for auth_user_id: %s", userID)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve seeker", "details": err.Error()})
			log.Printf("Error retrieving seeker for auth_user_id: %s, Error: %s", userID, err.Error())
		}
		return
	}

	isFilled := handlers.IsFieldFilled(seeker.ProfessionalSummary)
	if !isFilled {
		c.JSON(http.StatusNotFound, gin.H{"error": "Professional summary not yet filled"})
		return
	}

	profSummary, err := handlers.GetProfessionalSummary(&seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse professional summary", "details": err.Error()})
		log.Printf("Failed to parse professional summary for auth_user_id: %s, Error: %s", userID, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.ProfessionalSummaryResponse{
		AuthUserID:   userID,
		About:        profSummary.About,
		Skills:       profSummary.Skills,
		AnnualIncome: profSummary.AnnualIncome,
	})
}


func (h *ProfessionalSummaryHandler) UpdateProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	seekersCollection := c.MustGet("db").(*mongo.Database).Collection("seekers")

	var input dto.ProfessionalSummaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve seeker", "details": err.Error()})
		}
		return
	}

	if err := handlers.SetProfessionalSummary(&seeker, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update professional summary", "details": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"professional_summary": seeker.ProfessionalSummary,
		},
	}

	_, err = seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ProfessionalSummaryResponse{
		AuthUserID:   seeker.AuthUserID,
		About:        input.About,
		Skills:       input.Skills,
		AnnualIncome: input.AnnualIncome,
	})
}


func (h *ProfessionalSummaryHandler) PatchProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	seekersCollection := c.MustGet("db").(*mongo.Database).Collection("seekers")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve seeker", "details": err.Error()})
		}
		return
	}

	// Unmarshal existing professional summary into a map
	var profSummaryMap map[string]interface{}
	if err := handlers.GetFieldFromBson(seeker.ProfessionalSummary, &profSummaryMap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal professional summary", "details": err.Error()})
		return
	}

	// Apply the partial updates
	for key, value := range updates {
		profSummaryMap[key] = value
	}

	// Marshal updated map back to BSON
	if err := handlers.SetFieldToBson(profSummaryMap, &seeker.ProfessionalSummary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated professional summary", "details": err.Error()})
		return
	}

	_, err = seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, bson.M{
		"$set": bson.M{"professional_summary": seeker.ProfessionalSummary},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Professional summary patched successfully"})
}
