package dataentry

// import (
// 	"RAAS/dto"
// 	"RAAS/models"
// 	"encoding/json"
// 	"fmt"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// 	"net/http"
// 	"strconv"
// 	"time"
// )

// type EducationHandler struct {
// 	DB *gorm.DB
// }

// func NewEducationHandler(db *gorm.DB) *EducationHandler {
// 	return &EducationHandler{DB: db}
// }

// func (h *EducationHandler) CreateEducation(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	var input dto.EducationRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
// 		return
// 	}

// 	// Validate that none of the required fields are empty
// 	if input.Degree == "" || input.Institution == "" || input.FieldOfStudy == "" || input.StartDate.IsZero() || input.Achievements == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
// 		return
// 	}

// 	// Fetch seeker from the database
// 	var seeker models.Seeker
// 	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	// Check if Education is empty or null, and initialize it as an empty slice if necessary
// 	var educations []map[string]interface{}
// 	if len(seeker.Educations) > 0 {
// 		if err := json.Unmarshal(seeker.Educations, &educations); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse educations", "details": err.Error()})
// 			return
// 		}
// 	}

// 	// Create a new education entry
// 	newEducation := map[string]interface{}{
// 		"degree":        input.Degree,
// 		"institution":   input.Institution,
// 		"fieldOfStudy":  input.FieldOfStudy,
// 		"startDate":     input.StartDate.Format("2006-01-02"),
// 		"endDate":       input.EndDate.Format("2006-01-02"),
// 		"achievements":  input.Achievements,
// 	}

// 	// Append the new education entry
// 	educations = append(educations, newEducation)

// 	// Convert the updated education data back to JSON
// 	updatedEducationsJSON, err := json.Marshal(educations)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated educations", "details": err.Error()})
// 		return
// 	}

// 	// Update the educations in the database
// 	seeker.Educations = updatedEducationsJSON
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update educations", "details": err.Error()})
// 		return
// 	}

// 	// Create a response object for the new education
// 	response := dto.EducationResponse{
// 		ID:             uint(len(educations)), // Dynamically generate ID
// 		AuthUserID:     userID,
// 		Degree:         input.Degree,
// 		Institution:    input.Institution,
// 		FieldOfStudy:   input.FieldOfStudy,
// 		StartDate:      input.StartDate,
// 		EndDate:        input.EndDate,
// 		Achievements:   input.Achievements,
// 	}

// 	// Return the response with the new education
// 	c.JSON(http.StatusCreated, response)
// }

// func (h *EducationHandler) GetEducations(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	var seeker models.Seeker
// 	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	if len(seeker.Educations) == 0 || string(seeker.Educations) == "null" {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "No education records found"})
// 		return
// 	}

// 	var educations []map[string]interface{}
// 	if err := json.Unmarshal(seeker.Educations, &educations); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse education records"})
// 		return
// 	}

// 	var response []dto.EducationResponse
// 	for idx, edu := range educations {
// 		startDate, err := time.Parse("2006-01-02", edu["startDate"].(string))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid start date format"})
// 			return
// 		}

// 		// Parse EndDate directly as it's now required
// 		endDate, err := time.Parse("2006-01-02", edu["endDate"].(string))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid end date format"})
// 			return
// 		}

// 		response = append(response, dto.EducationResponse{
// 			ID:             uint(idx + 1),
// 			AuthUserID:     userID,
// 			Degree:         edu["degree"].(string),
// 			Institution:    edu["institution"].(string),
// 			FieldOfStudy:   edu["fieldOfStudy"].(string),
// 			StartDate:      startDate,
// 			EndDate:        endDate,
// 			Achievements:   edu["achievements"].(string),
// 		})
// 	}

// 	c.JSON(http.StatusOK, response)
// }

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
