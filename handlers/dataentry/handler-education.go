package dataentry

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EducationHandler struct
type EducationHandler struct {
	DB *gorm.DB
}

// NewEducationHandler creates a new EducationHandler
func NewEducationHandler(db *gorm.DB) *EducationHandler {
	return &EducationHandler{DB: db}
}

// CreateEducation creates a new education record for the authenticated user
func (h *EducationHandler) CreateEducation(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input dto.EducationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	education := models.Education{
		AuthUserID:   userID,
		Degree:       input.Degree,
		Institution:  input.Institution,
		FieldOfStudy: input.FieldOfStudy,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		Achievements: input.Achievements,
	}

	if err := h.DB.Create(&education).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create education", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.EducationResponse{
		ID:           education.ID,
		AuthUserID:   userID,
		Degree:       education.Degree,
		Institution:  education.Institution,
		FieldOfStudy: education.FieldOfStudy,
		StartDate:    education.StartDate,
		EndDate:      education.EndDate,
		Achievements: education.Achievements,
	})
}

// GetEducation retrieves all education records for the authenticated user
func (h *EducationHandler) GetEducation(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var records []models.Education
	if err := h.DB.Where("auth_user_id = ?", userID).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch education records", "details": err.Error()})
		return
	}

	var response []dto.EducationResponse
	for _, ed := range records {
		response = append(response, dto.EducationResponse{
			ID:           ed.ID,
			AuthUserID:   ed.AuthUserID,
			Degree:       ed.Degree,
			Institution:  ed.Institution,
			FieldOfStudy: ed.FieldOfStudy,
			StartDate:    ed.StartDate,
			EndDate:      ed.EndDate,
			Achievements: ed.Achievements,
		})
	}

	c.JSON(http.StatusOK, response)
}

// PutEducation updates an existing education record for the authenticated user
func (h *EducationHandler) PutEducation(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	var existing models.Education
	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Education record not found"})
		return
	}

	var updated models.Education
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Ensure these values stay correct regardless of input
	updated.ID = existing.ID
	updated.AuthUserID = userID
	updated.CreatedAt = existing.CreatedAt

	if err := h.DB.Save(&updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update education", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Education updated"})
}

// DeleteEducation deletes an existing education record for the authenticated user
func (h *EducationHandler) DeleteEducation(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id := c.Param("id")

	if err := h.DB.Where("id = ? AND auth_user_id = ?", id, userID).Delete(&models.Education{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete education", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Education deleted"})
}
