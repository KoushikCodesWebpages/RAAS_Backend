package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"


	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"RAAS/models"
	"RAAS/config"


	"RAAS/handlers/repo"
	// "bytes"
	// "log"
)

type CVData struct {
	JobID            string                 `json:"job_id"`
	Name             string                 `json:"name"`
	Designation      string                 `json:"designation"`
	Contact          string                 `json:"contact"`
	ProfileSummary   string                 `json:"profile_summary"`
	SkillsAndTools   []string               `json:"skills_and_tools"`
	Education        []EducationData        `json:"education"`
	ExperienceSummary []ExperienceSummaryData `json:"experience_summary"`
	Languages        []string               `json:"languages"`
}


type EducationData struct {
	Years       string   `json:"years"`
	Institution string   `json:"institution"`
	Details     []string `json:"details"`
}

type ExperienceSummaryData struct {
	Title  string   `json:"title"`
	Bullets []string `json:"bullets"`
}



// CVRequest struct to receive CV generation request
type CVRequest struct {
	// Add any required fields here
}

// CVHandler struct
type CVHandler struct {
	db     *gorm.DB
}

func NewCVHandler(db *gorm.DB, cfg *config.Config) *CVHandler {
	return &CVHandler{
		db:     db,
	}
}
func (h *CVHandler) PostCV(c *gin.Context) {

	var req CoverLetterAndCVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid job_id in request body"})
		return
	}
	jobID := req.JobID
	// Step 1: Extract user information from JWT claims
	userID := c.MustGet("userID").(uuid.UUID)

	// Step 2: Get user personal info

	//NAME
	var personalInfo models.PersonalInfo
	if err := h.db.Where("auth_user_id = ?", userID).First(&personalInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user personal information"})
		return
	}

	var fullName string
	if personalInfo.SecondName != nil {
		fullName = fmt.Sprintf("%s %s", personalInfo.FirstName, *personalInfo.SecondName)
	} else {
		fullName = personalInfo.FirstName
	}

	//ADDRESS
	address := personalInfo.Address

	//LINK
	linkedinlink:=personalInfo.LinkedInProfile


	//Step 3: Get seeker model

	var seeker models.Seeker
	if err := h.db.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve seeker data"})
		return
	}

	if seeker.DailyGeneratableCV > 0 {
		seeker.DailyGeneratableCV -= 1
		if err := h.db.Save(&seeker).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker data"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Daily CV generation limit exceeded"})
		return
	}

	//Step 4: Get Contact details.

	var authUser models.AuthUser
	if err := h.db.Where("id = ?", userID).First(&authUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve auth user info"})
		return
	}

	//Email PHONE
	email := authUser.Email
	phone := authUser.Phone

	// Step 5: Get professional summary (includes skills)
	var summary models.ProfessionalSummary
	if err := h.db.Where("auth_user_id = ?", userID).First(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve professional summary"})
		return
	}
	profile_summary := summary.About


	var skills []string
	if err := json.Unmarshal(summary.Skills, &skills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse skills"})
		return
	}

	// Step 6: Get education

	var educations []models.Education
	if err := h.db.Where("auth_user_id = ?", userID).Find(&educations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve education details"})
		return
	}
	
	educationData := make([]struct {
		Years      string   `json:"years"`
		Institution string `json:"institution"`
		Details    []string `json:"details"`
	}, 0)
	for _, edu := range educations {
		startYear := edu.StartDate.Format("2006")
		endYear := "Present"
		if edu.EndDate != nil {
			endYear = edu.EndDate.Format("2006")
		}
		details := make([]string, 0)
		details = append(details, edu.Degree)
		if edu.Achievements != "" {
			details = append(details, edu.Achievements)
		}
		educationData = append(educationData, struct {
			Years      string   `json:"years"`
			Institution string `json:"institution"`
			Details    []string `json:"details"`
		}{
			Years: fmt.Sprintf("%s - %s", startYear, endYear),
			Institution: edu.Institution,
			Details: details,
		})
	}

	// Step 7: Get work experience
	var workExperiences []models.WorkExperience
	if err := h.db.Where("auth_user_id = ?", userID).Find(&workExperiences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve work experience"})
		return
	}
	var workExpText string
	for _, we := range workExperiences {
		workExpText += fmt.Sprintf("Job Title: %s at %s. Responsibilities: %s. ", we.JobTitle, we.CompanyName, we.KeyResponsibilities)
	}

	experienceSummaryData := make([]struct {
		Title  string   `json:"title"`
		Bullets []string `json:"bullets"`
	}, 0)
	for _, we := range workExperiences {
		startYear := we.StartDate.Format("2006")
		endYear := "PRESENT"
		if we.EndDate != nil {
			endYear = we.EndDate.Format("2006")
		}
		bullets := strings.Split(we.KeyResponsibilities, ".")
		for i, bullet := range bullets {
			bullets[i] = strings.TrimSpace(bullet)
		}
		bullets = removeEmptyStrings(bullets)
		experienceSummaryData = append(experienceSummaryData, struct {
			Title  string   `json:"title"`
			Bullets []string `json:"bullets"`
		}{
			Title: fmt.Sprintf("%s at %s (%s - %s)", we.JobTitle, we.CompanyName, startYear, endYear),
			Bullets: bullets,
		})
	}

	// Step 8: Get Language

	var languages []models.Language
	if err := h.db.Where("auth_user_id = ?", userID).Find(&languages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve languages"})
		return
	}

	languageData := make([]string, 0)
	for _, lang := range languages {
		languageData = append(languageData, fmt.Sprintf("%s - %s", lang.LanguageName, lang.ProficiencyLevel))
	}

	var linkedin string
	if linkedinlink != nil {
		linkedin = *linkedinlink
	}

	// Step 9 Get Job Title

	var jobTitle string
	var jobMetaData models.LinkedInJobMetaData
	if err := h.db.Where("job_id = ?", jobID).First(&jobMetaData).Error; err == nil {
		jobTitle = jobMetaData.Title
	} else {
		var xingJobMetaData models.XingJobMetaData
		if err := h.db.Where("job_id = ?", jobID).First(&xingJobMetaData).Error; err == nil {
			jobTitle = xingJobMetaData.Title
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job metadata"})
			return
		}
	}

	contact := fmt.Sprintf("Email: %s\nPhone: %s\nAddress: %s\nLinkedIn: %s", email, phone, address, linkedin)

	cvInput := repo.CVInput{
		Name:             fullName,
		Designation:      jobTitle,
		Contact:          contact,
		ProfileSummary:   profile_summary,
		SkillsAndTools:   skills,
		Education:        educationData,
		ExperienceSummary: experienceSummaryData,
		Languages:        languageData,
	}

	cvData, err := repo.GenerateCVDocx(cvInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CV document", "details": err.Error()})
		return
	}
	result := h.db.Model(&models.SelectedJobApplication{}).Where("auth_user_id = ? AND job_id = ?", userID, jobID).Update("cv_generated", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update CvGenerated field"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=cv.docx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", cvData)
}


	// Step 7: Generate CV using external service or template



func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}


