package repo

import (
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type SeekerDTO struct {
	ID                        uint           `json:"id"`
	AuthUserID                uuid.UUID      `json:"authUserId"`
	SubscriptionTier          string         `json:"subscriptionTier"`
	DailySelectableJobsCount  int            `json:"dailySelectableJobsCount"`
	DailyGeneratableCV        int            `json:"dailyGeneratableCv"`
	DailyGeneratableCoverletter int          `json:"dailyGeneratableCoverletter"`
	TotalApplications         int            `json:"totalApplications"`
	PersonalInfo              interface{}    `json:"personalInfo"`  // Storing as JSON (interface{} for flexibility)
	ProfessionalSummary      interface{}    `json:"professionalSummary"`
	WorkExperiences           interface{}    `json:"workExperiences"`
	Educations                interface{}    `json:"education"`
	Certificates              interface{}    `json:"certificates"`
	Languages                 interface{}    `json:"languages"`
	PrimaryTitle             string         `json:"primaryTitle"`
	SecondaryTitle           *string        `json:"secondaryTitle,omitempty"`
	TertiaryTitle            *string        `json:"tertiaryTitle,omitempty"`
}

// SeekerHandler struct for managing seeker data
type SeekerHandler struct {
	DB *gorm.DB
}

// NewSeekerHandler returns a new SeekerHandler instance
func NewSeekerHandler(db *gorm.DB) *SeekerHandler {
	return &SeekerHandler{DB: db}
}

// GetSeeker - Get a seeker profile data based on AuthUserID
func (h *SeekerHandler) GetSeeker(c *gin.Context) {
	// Get user ID from JWT claims
	userID := c.MustGet("userID").(uuid.UUID)

	// Fetch Seeker data
	var seeker models.Seeker
	if err := h.DB.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Seeker not found",
		})
		return
	}

	// Prepare response with Seeker data
	seekerDTO := SeekerDTO{
		ID:                      seeker.ID,
		AuthUserID:              seeker.AuthUserID,
		SubscriptionTier:        seeker.SubscriptionTier,
		DailySelectableJobsCount: seeker.DailySelectableJobsCount,
		DailyGeneratableCV:      seeker.DailyGeneratableCV,
		DailyGeneratableCoverletter: seeker.DailyGeneratableCoverletter,
		TotalApplications:       seeker.TotalApplications,
		PersonalInfo:            seeker.PersonalInfo,
		ProfessionalSummary:    seeker.ProfessionalSummary,
		WorkExperiences:        seeker.WorkExperiences,
		Educations:             seeker.Educations,
		Certificates:           seeker.Certificates,
		Languages:              seeker.Languages,
		PrimaryTitle:           seeker.PrimaryTitle,
		SecondaryTitle:         seeker.SecondaryTitle,
		TertiaryTitle:          seeker.TertiaryTitle,
	}

	// Respond with the Seeker data
	c.JSON(http.StatusOK, gin.H{
		"seeker": seekerDTO,
	})
}
