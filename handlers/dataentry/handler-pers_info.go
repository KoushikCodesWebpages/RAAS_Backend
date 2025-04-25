package dataentry

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"RAAS/dto"
	"RAAS/models"
	"RAAS/handlers"
	"log"
)

type PersonalInfoHandler struct{}

func NewPersonalInfoHandler() *PersonalInfoHandler {
	return &PersonalInfoHandler{}
}

func (h *PersonalInfoHandler) CreatePersonalInfo(c *gin.Context) {
	// Get user ID from context
	userID := c.MustGet("userID").(uuid.UUID).String()
	// Get MongoDB collection from context
	collection := c.MustGet("db").(*mongo.Database).Collection("seekers")

	// Bind input data to PersonalInfoRequest
	var input dto.PersonalInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		log.Printf("Error binding input: %s", err.Error()) // Log binding error
		return
	}

	// Log the received input for debugging
	log.Printf("Received input: %+v", input)

	// Ensure DateOfBirth is not empty
	if input.DateOfBirth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date of birth cannot be empty"})
		log.Printf("Date of birth is empty for userID: %s", userID) // Log missing date of birth
		return
	}

	// Set context and timeout for MongoDB operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	// Fetch seeker by "auth_user_id"
	err := collection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
			log.Printf("Seeker not found for userID: %s", userID) // Log if seeker not found
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving seeker"})
			log.Printf("Error retrieving seeker for userID: %s, Error: %s", userID, err.Error()) // Log DB error
		}
		return
	}

	// Log the seeker document for debugging
	log.Printf("Found seeker: %+v", seeker)

	// Check if personal info is already set
	if handlers.IsFieldFilled(seeker.PersonalInfo) {
		c.JSON(http.StatusConflict, gin.H{"error": "Personal info already exists"})
		log.Printf("Personal info already exists for userID: %s", userID) // Log conflict
		return
	}

	// Set the personal info from the request
	if err := handlers.SetPersonalInfo(&seeker, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process personal info"})
		log.Printf("Failed to process personal info for userID: %s, Error: %s", userID, err.Error()) // Log processing error
		return
	}

	// Log the personal info after setting it
	log.Printf("Set personal info for userID: %s: %+v", userID, seeker.PersonalInfo)

	// Update personal info in MongoDB
	updateResult, err := collection.UpdateOne(ctx, bson.M{"auth_user_id": userID}, bson.M{"$set": bson.M{"personal_info": seeker.PersonalInfo}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save personal info"})
		log.Printf("Failed to update personal info for userID: %s, Error: %s", userID, err.Error()) // Log update error
		return
	}

	// Log the result of the update operation
	log.Printf("Update result: %+v", updateResult)

	// Check if the update operation affected any documents
	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching seeker found to update"})
		log.Printf("No matching seeker found for userID: %s", userID) // Log no matching seeker
		return
	}

	// Return the newly created personal info in the response
	c.JSON(http.StatusCreated, dto.PersonalInfoResponse{
		AuthUserID:      userID,
		FirstName:       input.FirstName,
		SecondName:      input.SecondName,
		DateOfBirth:     input.DateOfBirth,
		Address:         input.Address,
		LinkedInProfile: input.LinkedInProfile,
	})

	// Log the response being returned
	log.Printf("Returning response: %+v", dto.PersonalInfoResponse{
		AuthUserID:      userID,
		FirstName:       input.FirstName,
		SecondName:      input.SecondName,
		DateOfBirth:     input.DateOfBirth,
		Address:         input.Address,
		LinkedInProfile: input.LinkedInProfile,
	})
}



func (h *PersonalInfoHandler) GetPersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID).String()
	collection := c.MustGet("db").(*mongo.Database).Collection("seekers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&seeker)
	if err != nil || seeker.PersonalInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personal info not found"})
		return
	}

	personalInfo, err := handlers.GetPersonalInfo(&seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal personal info"})
		return
	}

	c.JSON(http.StatusOK, personalInfo)
}

func (h *PersonalInfoHandler) UpdatePersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID).String()
	collection := c.MustGet("db").(*mongo.Database).Collection("seekers")

	var input dto.PersonalInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&seeker)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	if err := handlers.SetPersonalInfo(&seeker, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal personal info"})
		return
	}

	_, err = collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"personal_info": seeker.PersonalInfo}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info"})
		return
	}

	c.JSON(http.StatusOK, dto.PersonalInfoResponse{
		AuthUserID:      userID,
		FirstName:       input.FirstName,
		SecondName:      input.SecondName,
		DateOfBirth:     input.DateOfBirth,
		Address:         input.Address,
		LinkedInProfile: input.LinkedInProfile,
	})
}


func (h *PersonalInfoHandler) PatchPersonalInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID).String()
	collection := c.MustGet("db").(*mongo.Database).Collection("seekers")

	var input dto.PersonalInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.DateOfBirth != "" || input.Address != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot patch DateOfBirth or Address"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var seeker models.Seeker
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&seeker)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	personalInfo, err := handlers.GetPersonalInfo(&seeker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve personal info"})
		return
	}

	// Patch allowed fields
	if input.FirstName != "" {
		personalInfo.FirstName = input.FirstName
	}
	if input.SecondName != nil {
		personalInfo.SecondName = input.SecondName
	}
	if input.LinkedInProfile != nil {
		personalInfo.LinkedInProfile = input.LinkedInProfile
	}

	if err := handlers.SetPersonalInfo(&seeker, personalInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info"})
		return
	}

	_, err = collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"personal_info": seeker.PersonalInfo}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save patched personal info"})
		return
	}

	c.JSON(http.StatusOK, personalInfo)
}
