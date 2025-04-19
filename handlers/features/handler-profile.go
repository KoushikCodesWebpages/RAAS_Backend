package features

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"net/http"
	"time"
	"encoding/json"
	"fmt"
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
func (h *SeekerProfileHandler) GetSeekerProfile(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)

    var seeker models.Seeker
    if err := h.DB.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
        c.JSON(http.StatusNoContent, gin.H{"error": "Seeker not found"})
        return
    }

    // Unmarshal PersonalInfo
    var personalInfo struct {
        FirstName  string  `json:"firstName"`
        SecondName *string `json:"secondName"`
    }
    if err := json.Unmarshal(seeker.PersonalInfo, &personalInfo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse personal info"})
        return
    }

    // Unmarshal ProfessionalSummary
    var summary struct {
        Skills []string `json:"skills"`
    }
    if err := json.Unmarshal(seeker.ProfessionalSummary, &summary); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse professional summary"})
        return
    }

    // Unmarshal WorkExperiences
    var workExperiences []map[string]interface{}
    if len(seeker.WorkExperiences) > 0 {
        if err := json.Unmarshal(seeker.WorkExperiences, &workExperiences); err != nil {
            fmt.Printf("Failed to unmarshal work experiences: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse work experiences", "details": err.Error()})
            return
        }
    }

    // Total experience in months
    totalExperienceInMonths := 0
    for _, w := range workExperiences {
        startDate, err := time.Parse("2006-01-02", w["startDate"].(string))
        if err != nil {
            fmt.Printf("Error parsing start date: %v\n", err)
            continue
        }

        var endDate time.Time
        if endStr, ok := w["endDate"].(string); ok && endStr != "" {
            endDate, err = time.Parse("2006-01-02", endStr)
            if err != nil {
                fmt.Printf("Error parsing end date: %v\n", err)
                continue
            }
        } else {
            endDate = time.Now() // Current time if end date is missing
        }

        duration := endDate.Sub(startDate)
        totalExperienceInMonths += int(duration.Hours() / 24 / 30) // Convert to months
    }

    // Unmarshal Certificates
    var certificates []struct {
        CertificateName string `json:"certificateName"`
    }
    if err := json.Unmarshal(seeker.Certificates, &certificates); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse certificates"})
        return
    }

    certNames := make([]string, 0, len(certificates))
    for _, cert := range certificates {
        certNames = append(certNames, cert.CertificateName)
    }

    // Unmarshal Languages
    var languages []struct {
        LanguageName string `json:"languageName"`
    }
    if err := json.Unmarshal(seeker.Languages, &languages); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse languages"})
        return
    }

    // Count total available jobs
    var totalJobs int64
    if err := h.DB.Model(&models.Job{}).Count(&totalJobs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching total jobs"})
        return
    }

    // Profile completion logic
    completion := 70
    if seeker.SubscriptionTier != "" &&
        personalInfo.FirstName != "" &&
        len(summary.Skills) > 0 &&
        totalExperienceInMonths > 0 &&
        len(certNames) > 0 &&
        len(languages) > 0 &&
        seeker.PrimaryTitle != "" {
        completion = 100
    }

    // Build and return DTO
    dto := dto.SeekerProfileDTO{
        ID:                          seeker.ID,
        AuthUserID:                  seeker.AuthUserID,
        FirstName:                   personalInfo.FirstName,
        SecondName:                  personalInfo.SecondName,
        Skills:                      summary.Skills,
        TotalExperienceInMonths:     totalExperienceInMonths,
        Certificates:                certNames,
        PreferredJobTitle:           seeker.PrimaryTitle,
        SubscriptionTier:            seeker.SubscriptionTier,
        DailySelectableJobsCount:    seeker.DailySelectableJobsCount,
        DailyGeneratableCV:          seeker.DailyGeneratableCV,
        DailyGeneratableCoverletter: seeker.DailyGeneratableCoverletter,
        TotalApplications:           seeker.TotalApplications,
        TotalJobsAvailable:          int(totalJobs),
        ProfileCompletion:           completion,
    }

    c.JSON(http.StatusOK, gin.H{"profile": dto})
}
