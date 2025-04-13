package dataentry

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorkExperienceHandler struct
type WorkExperienceHandler struct {
	DB *gorm.DB
}

// NewWorkExperienceHandler creates a new WorkExperienceHandler
func NewWorkExperienceHandler(db *gorm.DB) *WorkExperienceHandler {
	return &WorkExperienceHandler{DB: db}
}

// CreateWorkExperience creates a work experience for the authenticated user
func (h *WorkExperienceHandler) CreateWorkExperience(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.WorkExperienceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	workExp := models.WorkExperience{
		AuthUserID:         userID,
		JobTitle:           input.JobTitle,
		CompanyName:        input.CompanyName,
		EmployerType:       input.EmployerType,
		StartDate:          input.StartDate,
		EndDate:            input.EndDate,
		KeyResponsibilities: input.KeyResponsibilities,
	}

	if err := h.DB.Create(&workExp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create work experience", "details": err.Error()})
		return
	}

	var timeline models.UserEntryTimeline
    if err := h.DB.First(&timeline, "user_id = ?", userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
        return
    }

    timeline.WorkExperiencesCompleted = true

    // Save the updated timeline
    if err := h.DB.Save(&timeline).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timeline", "details": err.Error()})
        return
    }

	response := dto.WorkExperienceResponse{
		ID:                 workExp.ID,
		AuthUserID:         userID,
		JobTitle:           workExp.JobTitle,
		CompanyName:        workExp.CompanyName,
		EmployerType:       workExp.EmployerType,
		StartDate:          workExp.StartDate,
		EndDate:            workExp.EndDate,
		KeyResponsibilities: workExp.KeyResponsibilities,
	}

	c.JSON(http.StatusCreated, response)
}

// GetWorkExperience retrieves the work experiences of the authenticated user
func (h *WorkExperienceHandler) GetWorkExperience(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var workExps []models.WorkExperience
	if err := h.DB.Where("auth_user_id = ?", userID).Find(&workExps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch work experiences", "details": err.Error()})
		return
	}

	var response []dto.WorkExperienceResponse
	for _, we := range workExps {
		response = append(response, dto.WorkExperienceResponse{
			ID:                 we.ID,
			AuthUserID:         we.AuthUserID,
			JobTitle:           we.JobTitle,
			CompanyName:        we.CompanyName,
			EmployerType:       we.EmployerType,
			StartDate:          we.StartDate,
			EndDate:            we.EndDate,
			KeyResponsibilities: we.KeyResponsibilities,
		})
	}

	c.JSON(http.StatusOK, response)
}

// PatchWorkExperience partially updates the work experience of the authenticated user
func (h *WorkExperienceHandler) PatchWorkExperience(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var workExp models.WorkExperience
	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).First(&workExp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work experience not found"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if err := h.DB.Model(&workExp).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work experience", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work experience updated"})
}

// DeleteWorkExperience deletes the work experience of the authenticated user
func (h *WorkExperienceHandler) DeleteWorkExperience(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).Delete(&models.WorkExperience{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete work experience", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work experience deleted"})
}
