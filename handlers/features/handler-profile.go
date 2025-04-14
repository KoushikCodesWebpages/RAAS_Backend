package features

import (

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"net/http"
	"time"
	"encoding/json"

	"RAAS/models"
	"RAAS/dto"
)

// SeekerProfileHandler struct for managing seeker profile
type SeekerProfileHandler struct {
	DB *gorm.DB
}

// NewSeekerProfileHandler returns a new SeekerProfileHandler instance
func NewSeekerProfileHandler(db *gorm.DB) *SeekerProfileHandler {
	return &SeekerProfileHandler{DB: db}
}

// GetSeekerProfile - Get all the data related to a user's profile from multiple sources
func (h *SeekerProfileHandler) GetSeekerProfile(c *gin.Context) {
	// Get user ID from JWT claims
	userID := c.MustGet("userID").(uuid.UUID)

	// Fetch Seeker data
	var seeker models.Seeker
	if err := h.DB.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Seeker not found",
		})
		return
	}

	// Fetch PersonalInfo data
	var personalInfo models.PersonalInfo
	if err := h.DB.Where("auth_user_id = ?", userID).First(&personalInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Personal info not found",
		})
		return
	}

	// Fetch ProfessionalSummary data
	var professionalSummary models.ProfessionalSummary
	if err := h.DB.Where("auth_user_id = ?", userID).First(&professionalSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Professional summary not found",
		})
		return
	}

	// Fetch WorkExperience data and calculate total experience in months
	var workExperiences []models.WorkExperience
	if err := h.DB.Where("auth_user_id = ?", userID).Find(&workExperiences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Work experience not found",
		})
		return
	}
	totalExperienceInMonths := 0
	for _, work := range workExperiences {
		// Calculate months of experience for each work entry
		var endDate time.Time
		if work.EndDate != nil {
			endDate = *work.EndDate
		} else {
			endDate = time.Now()
		}
		experienceDuration := endDate.Sub(work.StartDate)
		totalExperienceInMonths += int(experienceDuration.Hours() / 24 / 30) // Rough estimate in months
	}

	// Fetch Certificate data
	var certificates []models.Certificate
	if err := h.DB.Where("auth_user_id = ?", userID).Find(&certificates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Certificates not found",
		})
		return
	}
	certificateNames := []string{}
	for _, cert := range certificates {
		certificateNames = append(certificateNames, cert.CertificateName)
	}

	// Fetch Language data
	var languages []models.Language
	if err := h.DB.Where("auth_user_id = ?", userID).Find(&languages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Languages not found",
		})
		return
	}
	languageNames := []string{}
	for _, lang := range languages {
		languageNames = append(languageNames, lang.LanguageName)
	}

	// Fetch PreferredJobTitle data
	var preferredJobTitle models.PreferredJobTitle
	if err := h.DB.Where("auth_user_id = ?", userID).First(&preferredJobTitle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Preferred job title not found",
		})
		return
	}

	// Fetch total number of available jobs from all sources
	var totalJobs int
	var linkedInJobs []models.LinkedInJobMetaData
	var xingJobs []models.XingJobMetaData
	if err := h.DB.Find(&linkedInJobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching LinkedIn job data",
		})
		return
	}
	if err := h.DB.Find(&xingJobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching Xing job data",
		})
		return
	}
	totalJobs = len(linkedInJobs) + len(xingJobs)

	// Calculate profile completion percentage
	completionPercentage := 0
	// Assume all required fields must be non-empty for completion
	if seeker.SubscriptionTier != "" && personalInfo.FirstName != "" && professionalSummary.Skills != nil && totalExperienceInMonths > 0 && len(certificates) > 0 && len(languageNames) > 0 && preferredJobTitle.PrimaryTitle != "" {
		completionPercentage = 100
	} else {
		completionPercentage = 70 // Example value, you could calculate based on which fields are present
	}

	var skills []string
	if err := json.Unmarshal(professionalSummary.Skills, &skills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse skills",
		})
		return
	}


	// Create SeekerProfileDTO response
	profileDTO := dto.SeekerProfileDTO{
		ID:                       seeker.ID,
		AuthUserID:               seeker.AuthUserID,
		FirstName:                personalInfo.FirstName,
		SecondName:               personalInfo.SecondName,
		Skills: skills,
		TotalExperienceInMonths:  totalExperienceInMonths,
		Certificates:             certificateNames,
		PreferredJobTitle:       preferredJobTitle.PrimaryTitle,
		SubscriptionTier:         seeker.SubscriptionTier,
		DailySelectableJobsCount: seeker.DailySelectableJobsCount,
		DailyGeneratableCV:       seeker.DailyGeneratableCV,
		DailyGeneratableCoverletter: seeker.DailyGeneratableCoverletter,
		TotalApplications:        seeker.TotalApplications,
		TotalJobsAvailable:       totalJobs,
		ProfileCompletion:        completionPercentage,
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": profileDTO,
	})
}
