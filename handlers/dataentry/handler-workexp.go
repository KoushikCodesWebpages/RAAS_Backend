package dataentry

import (
	"RAAS/dto"
	"RAAS/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type WorkExperienceHandler struct {
	DB *gorm.DB
}

func NewWorkExperienceHandler(db *gorm.DB) *WorkExperienceHandler {
	return &WorkExperienceHandler{DB: db}
}
func (h *WorkExperienceHandler) CreateWorkExperience(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)

    var input dto.WorkExperienceRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    // Validate that none of the required fields are empty
    if input.JobTitle == "" || input.CompanyName == "" || input.EmploymentType == "" || input.StartDate.IsZero() || input.KeyResponsibilities == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
        return
    }

    // Fetch seeker from the database
    var seeker models.Seeker
    if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
        return
    }

    // Check if WorkExperiences is empty or null, and initialize it as an empty slice if necessary
    var workExperiences []map[string]interface{}
    if len(seeker.WorkExperiences) > 0 {
        if err := json.Unmarshal(seeker.WorkExperiences, &workExperiences); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse work experiences", "details": err.Error()})
            return
        }
    }

    // Create a new work experience entry
    newWorkExperience := map[string]interface{}{
        "jobTitle":          input.JobTitle,
        "companyName":       input.CompanyName,
        "employmentType":    input.EmploymentType,
        "startDate":         input.StartDate.Format("2006-01-02"),
        "endDate":           input.EndDate.Format("2006-01-02"),
        "keyResponsibilities": input.KeyResponsibilities,
    }

    // Append the new work experience entry
    workExperiences = append(workExperiences, newWorkExperience)

    // Convert the updated work experiences back to JSON
    updatedWorkExperiencesJSON, err := json.Marshal(workExperiences)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated work experiences", "details": err.Error()})
        return
    }

    // Update the work experiences in the database
    seeker.WorkExperiences = updatedWorkExperiencesJSON
    if err := h.DB.Save(&seeker).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work experiences", "details": err.Error()})
        return
    }

    // Create a response object for the new work experience
    response := dto.WorkExperienceResponse{
        ID:                  uint(len(workExperiences)),  // Dynamically generate ID
        AuthUserID:          userID,
        JobTitle:            input.JobTitle,
        CompanyName:         input.CompanyName,
        EmploymentType:      input.EmploymentType,
        StartDate:           input.StartDate,
        EndDate:             input.EndDate,
        KeyResponsibilities: input.KeyResponsibilities,
    }

    // Return the response with the new work experience
    c.JSON(http.StatusCreated, response)
}


func (h *WorkExperienceHandler) GetWorkExperiences(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	if len(seeker.WorkExperiences) == 0 || string(seeker.WorkExperiences) == "null" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No work experiences found"})
		return
	}

	var workExperiences []map[string]interface{}
	if err := json.Unmarshal(seeker.WorkExperiences, &workExperiences); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse work experiences"})
		return
	}

	var response []dto.WorkExperienceResponse
	for idx, we := range workExperiences {
		startDate, err := time.Parse("2006-01-02", we["startDate"].(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid start date format"})
			return
		}

		// Parse EndDate directly as it's now required
		endDate, err := time.Parse("2006-01-02", we["endDate"].(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid end date format"})
			return
		}

		response = append(response, dto.WorkExperienceResponse{
			ID:                  uint(idx + 1),
			AuthUserID:          userID,
			JobTitle:            we["jobTitle"].(string),
			CompanyName:         we["companyName"].(string),
			EmploymentType:      we["employmentType"].(string),
			StartDate:           startDate,
			EndDate:             endDate,
			KeyResponsibilities: we["keyResponsibilities"].(string),
		})
	}

	c.JSON(http.StatusOK, response)
}


func (h *WorkExperienceHandler) PatchWorkExperience(c *gin.Context) {
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

	var workExperiences []map[string]interface{}
	if err := json.Unmarshal(seeker.WorkExperiences, &workExperiences); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse work experiences"})
		return
	}

	index, err := strconv.Atoi(id)
	if err != nil || index <= 0 || index > len(workExperiences) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work experience index"})
		return
	}

	// Apply updates
	entry := workExperiences[index-1]
	for key, value := range updateFields {
		if _, exists := entry[key]; exists {
			entry[key] = value
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid field: %s", key)})
			return
		}
	}
	workExperiences[index-1] = entry

	updatedJSON, err := json.Marshal(workExperiences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated work experiences"})
		return
	}

	seeker.WorkExperiences = updatedJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
		return
	}

	// Parse StartDate and EndDate (ensure both are valid)
	startDate, err := time.Parse("2006-01-02", entry["startDate"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid start date format"})
		return
	}

	// EndDate is required, so parse it directly without nil checks
	endDate, err := time.Parse("2006-01-02", entry["endDate"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid end date format"})
		return
	}

	// Create and return response with updated data
	response := dto.WorkExperienceResponse{
		ID:                  uint(index),
		AuthUserID:          userID,
		JobTitle:            entry["jobTitle"].(string),
		CompanyName:         entry["companyName"].(string),
		EmploymentType:      entry["employmentType"].(string),
		StartDate:           startDate,
		EndDate:             endDate,
		KeyResponsibilities: entry["keyResponsibilities"].(string),
	}

	c.JSON(http.StatusOK, response)
}

func (h *WorkExperienceHandler) DeleteWorkExperience(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var seeker models.Seeker
	if err := h.DB.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
		return
	}

	var workExperiences []map[string]interface{}
	if err := json.Unmarshal(seeker.WorkExperiences, &workExperiences); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse work experiences"})
		return
	}

	index, err := strconv.Atoi(id)
	if err != nil || index <= 0 || index > len(workExperiences) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work experience index"})
		return
	}

	// Remove the work experience at the specified index (index - 1 since it's 1-based in API)
	workExperiences = append(workExperiences[:index-1], workExperiences[index:]...)

	updatedJSON, err := json.Marshal(workExperiences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated work experiences"})
		return
	}

	seeker.WorkExperiences = updatedJSON
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work experience deleted successfully"})
}



