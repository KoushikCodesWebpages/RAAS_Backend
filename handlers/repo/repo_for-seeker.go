package repo

import (
	"RAAS/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv" 
)

// SeekerDTO is the Data Transfer Object for Seeker
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

	// Fetch AuthUser data associated with the seeker
	var authUser models.AuthUser
	if err := h.DB.Where("id = ?", seeker.AuthUserID).First(&authUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Auth user not found",
		})
		return
	}

	// Prepare response with Seeker and AuthUser data
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

	// Add AuthUser details to the response
	authUserDTO := struct {
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		CreatedAt string `json:"createdAt"`
	}{
		Email:     authUser.Email,    // Assuming there's an Email field in AuthUser model
		Phone:     authUser.Phone,    // Assuming there's a Phone field in AuthUser model
		CreatedAt: authUser.CreatedAt.String(), // Assuming CreatedAt is a time.Time
	}

	// Respond with Seeker and AuthUser data
	c.JSON(http.StatusOK, gin.H{
		"seeker":   seekerDTO,
		"authUser": authUserDTO,
	})
}

// GetAllSeekers - Get all seekers along with their AuthUser details
func (h *SeekerHandler) GetAllSeekers(c *gin.Context) {
	// Fetch all Seekers
	var seekers []models.Seeker
	if err := h.DB.Find(&seekers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch seekers",
		})
		return
	}

	// Prepare a slice to hold the response data
	var seekersWithAuth []gin.H

	// Loop through each seeker and fetch associated AuthUser data
	for _, seeker := range seekers {
		// Fetch AuthUser data associated with the seeker
		var authUser models.AuthUser
		if err := h.DB.Where("id = ?", seeker.AuthUserID).First(&authUser).Error; err != nil {
			// Handle error if AuthUser is not found
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Auth user not found for Seeker ID " + strconv.Itoa(int(seeker.ID)), // Convert uint to string
			})
			return
		}

		// Prepare response for each seeker with their auth details
		authUserDTO := struct {
			Email     string `json:"email"`
			Phone     string `json:"phone"`
			CreatedAt string `json:"createdAt"`
		}{
			Email:     authUser.Email,    // Assuming there's an Email field in AuthUser model
			Phone:     authUser.Phone,    // Assuming there's a Phone field in AuthUser model
			CreatedAt: authUser.CreatedAt.String(), // Assuming CreatedAt is a time.Time
		}

		// Prepare Seeker data
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

		// Add seeker and their auth details to the response slice
		seekersWithAuth = append(seekersWithAuth, gin.H{
			"seeker":   seekerDTO,
			"authUser": authUserDTO,
		})
	}

	// Respond with the list of seekers and their auth data
	c.JSON(http.StatusOK, gin.H{
		"seekers": seekersWithAuth,
	})
}