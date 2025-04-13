package features

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"fmt" // For logging
	"math/rand"
	"encoding/json" // For handling JSON unmarshalling
)

// JobRetrievalHandler handles the retrieval of jobs based on user's preferences and skills
func JobRetrievalHandler(c *gin.Context) {
	// Retrieve database and userID from context
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uuid.UUID)

	var jobs []dto.JobDTO

	// Fetch user's preferred job titles
	var preferred models.PreferredJobTitle
	if err := db.Where("auth_user_id = ?", userID).First(&preferred).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Set Your Job Title First",
		})
		return
	}

	// Collect preferred titles
	var preferredTitles []string
	if preferred.PrimaryTitle != "" {
		preferredTitles = append(preferredTitles, preferred.PrimaryTitle)
	}
	if preferred.SecondaryTitle != nil && *preferred.SecondaryTitle != "" {
		preferredTitles = append(preferredTitles, *preferred.SecondaryTitle)
	}
	if preferred.TertiaryTitle != nil && *preferred.TertiaryTitle != "" {
		preferredTitles = append(preferredTitles, *preferred.TertiaryTitle)
	}

	if len(preferredTitles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No preferred job titles set for user."})
		return
	}

	// Title filtering (case-insensitive)
	var conditions []string
	var values []interface{}
	for _, title := range preferredTitles {
		conditions = append(conditions, "LOWER(title) LIKE ?")
		values = append(values, "%"+strings.ToLower(title)+"%")
	}
	whereClause := strings.Join(conditions, " OR ")

	// Query jobs from both sources
	var linkedinJobs []models.LinkedInJobMetaData
	var xingJobs []models.XingJobMetaData

	// Fetch LinkedIn jobs
	db.Debug().Where(whereClause, values...).Find(&linkedinJobs)
	// Fetch Xing jobs
	db.Debug().Where(whereClause, values...).Find(&xingJobs)

	// Fetch user's professional summary to get skills
	var professionalSummary models.ProfessionalSummary
	if err := db.Where("auth_user_id = ?", userID).First(&professionalSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching professional summary",
		})
		return
	}

	// Get the user's skills from the professional summary (assuming it's a JSON array)
	var userSkills []string
	if professionalSummary.Skills != nil {
		// Unmarshal the JSON array into a Go slice
		if err := json.Unmarshal(professionalSummary.Skills, &userSkills); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error unmarshalling user skills",
			})
			return
		}
	}

	// Helper function to generate a random salary range (between 25k and 50k)
	randomSalary := func() dto.SalaryRange {
		minSalary := (rand.Intn(25) + 25) * 1000 // Random min salary between 25k and 50k
		maxSalary := (rand.Intn(25) + 25) * 1000 // Random max salary between 25k and 50k
		if minSalary > maxSalary {
			minSalary, maxSalary = maxSalary, minSalary // Ensure min is always less than max
		}
		return dto.SalaryRange{Min: minSalary, Max: maxSalary}
	}

	// LinkedIn Jobs Processing
	for _, job := range linkedinJobs {
		var jobDesc models.LinkedInJobDescription
		if err := db.Where("job_id = ?", job.JobID).First(&jobDesc).Error; err != nil {
			// Log the error when fetching job descriptions
			fmt.Println("Error fetching LinkedIn job description for job_id:", job.JobID, "Error:", err)
			continue
		}

		// Generate random salary range
		salaryRange := randomSalary()

		// Add job to jobs list
		jobs = append(jobs, dto.JobDTO{
			Source:         "linkedin", // Changed to 'linkedin' instead of 'xing' for LinkedIn jobs
			ID:             job.ID,
			JobID:          job.JobID,
			Title:          job.Title,
			Company:        job.Company,
			Location:       job.Location,
			PostedDate:     job.PostedDate,
			Processed:      job.Processed,
			JobType:        jobDesc.JobType,
			Skills:         jobDesc.Skills,
			UserSkills:     userSkills,
			ExpectedSalary: salaryRange,
			MatchScore:     0, // Initially 0, can be calculated later
			Description:    jobDesc.JobDescription, // Added line
		})
	}

	// Xing Jobs Processing
	for _, job := range xingJobs {
		var jobDesc models.XingJobDescription
		if err := db.Where("job_id = ?", job.JobID).First(&jobDesc).Error; err != nil {
			// Log the error when fetching job descriptions
			fmt.Println("Error fetching Xing job description for job_id:", job.JobID, "Error:", err)
			continue
		}

		// Generate random salary range
		salaryRange := randomSalary()

		// Add job to jobs list
		jobs = append(jobs, dto.JobDTO{
			Source:         "xing", // Corrected to 'xing' for Xing jobs
			ID:             job.ID,
			JobID:          job.JobID,
			Title:          job.Title,
			Company:        job.Company,
			Location:       job.Location,
			PostedDate:     job.PostedDate,
			Processed:      job.Processed,
			JobType:        jobDesc.JobType,
			Skills:         jobDesc.Skills,
			UserSkills:     userSkills,
			ExpectedSalary: salaryRange,
			MatchScore:     0, // Initially 0, can be calculated later
			Description:    jobDesc.JobDescription, // Added line
		})
	}

	// Respond with the jobs and pagination
	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
		"pagination": gin.H{
			"total": len(jobs),
		},
	})
}
