package dataentry

import (

	"RAAS/internal/dto"
	"RAAS/internal/handlers"
	"RAAS/internal/models"


	"context"
	"log"
	"net/http"
	"time"
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
	err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		handleDBError(err, c, "Seeker not found", userID)
		return
	}

	// Process and set professional summary
	if err := handlers.SetProfessionalSummary(&seeker, &input); err != nil {
		handleProcessingError(err, c, "Failed to process professional summary", userID)
		return
	}

	updateResult, err := seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, bson.M{"$set": bson.M{"professional_summary": seeker.ProfessionalSummary}})
	if err != nil || updateResult.MatchedCount == 0 {
		handleDBError(err, c, "Failed to save professional summary", userID)
		return
	}

	message := "Professional summary created"
	if handlers.IsFieldFilled(seeker.ProfessionalSummary) {
		message = "Professional summary updated"
	}

	// Update user entry timeline for the specific user
	timelineUpdate := bson.M{"$set": bson.M{"professional_summaries_completed": true}}
	timelineUpdateResult, err := entryTimelineCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, timelineUpdate)
	if err != nil || timelineUpdateResult.MatchedCount == 0 {
		handleDBError(err, c, "Failed to update user entry timeline", userID)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *ProfessionalSummaryHandler) GetProfessionalSummary(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	seekersCollection := c.MustGet("db").(*mongo.Database).Collection("seekers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		handleDBError(err, c, "Seeker not found", userID)
		return
	}

	if !handlers.IsFieldFilled(seeker.ProfessionalSummary) {
		c.JSON(http.StatusNoContent, gin.H{"error": "Professional summary not yet filled"})
		return
	}

	profSummary, err := handlers.GetProfessionalSummary(&seeker)
	if err != nil {
		handleProcessingError(err, c, "Failed to parse professional summary", userID)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"professional_summary": dto.ProfessionalSummaryResponse{
			AuthUserID:   userID,
			About:        profSummary.About,
			Skills:       profSummary.Skills,
			AnnualIncome: profSummary.AnnualIncome,
		},
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
		handleDBError(err, c, "Seeker not found", userID)
		return
	}

	if err := handlers.SetProfessionalSummary(&seeker, &input); err != nil {
		handleProcessingError(err, c, "Failed to update professional summary", userID)
		return
	}

	update := bson.M{"$set": bson.M{"professional_summary": seeker.ProfessionalSummary}}
	_, err = seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, update)
	if err != nil {
		handleDBError(err, c, "Database update failed", userID)
		return
	}

	c.JSON(http.StatusOK, dto.ProfessionalSummaryResponse{
		AuthUserID:   seeker.AuthUserID,
		About:        input.About,
		Skills:       input.Skills,
		AnnualIncome: input.AnnualIncome,
	})
}

// func (h *ProfessionalSummaryHandler) PatchProfessionalSummary(c *gin.Context) {
// 	userID := c.MustGet("userID").(string)
// 	seekersCollection := c.MustGet("db").(*mongo.Database).Collection("seekers")

// 	var updates map[string]interface{}
// 	if err := c.ShouldBindJSON(&updates); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	var seeker models.Seeker
// 	err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
// 	if err != nil {
// 		handleDBError(err, c, "Seeker not found", userID)
// 		return
// 	}

// 	var profSummaryMap map[string]interface{}
// 	if err := handlers.GetFieldFromBson(seeker.ProfessionalSummary, &profSummaryMap); err != nil {
// 		handleProcessingError(err, c, "Failed to unmarshal professional summary", userID)
// 		return
// 	}

// 	for key, value := range updates {
// 		profSummaryMap[key] = value
// 	}

// 	if err := handlers.SetFieldToBson(profSummaryMap, &seeker.ProfessionalSummary); err != nil {
// 		handleProcessingError(err, c, "Failed to marshal updated professional summary", userID)
// 		return
// 	}

// 	_, err = seekersCollection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, bson.M{
// 		"$set": bson.M{"professional_summary": seeker.ProfessionalSummary},
// 	})
// 	if err != nil {
// 		handleDBError(err, c, "Database update failed", userID)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Professional summary patched successfully"})
// }

// Helper function to handle DB errors
func handleDBError(err error, c *gin.Context, message, userID string) {
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": message})
		log.Printf("User %s: %s", userID, message)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": message, "details": err.Error()})
		log.Printf("Error for user %s: %s, Error: %s", userID, message, err.Error())
	}
}

// Helper function to handle processing errors
func handleProcessingError(err error, c *gin.Context, message, userID string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": message, "details": err.Error()})
	log.Printf("Processing error for user %s: %s, Error: %s", userID, message, err.Error())
}
