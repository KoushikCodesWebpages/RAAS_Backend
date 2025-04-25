
package features

// import (
// 	"RAAS/dto"
// 	"RAAS/models"
// 	"net/http"
// 	"strings"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// 	"fmt"
// 	"math/rand"
// 	"encoding/json"
// )

// // JobRetrievalHandler handles the retrieval of jobs based on user's preferences and skills
// // JobRetrievalHandler handles the retrieval of jobs based on user's preferences and skills
// func JobRetrievalHandler(c *gin.Context) {
// 	// Retrieve database and userID from context
// 	db := c.MustGet("db").(*gorm.DB)
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	var jobs []dto.JobDTO

// 	// Fetch user's preferred job titles from the Seeker model
// 	var seeker models.Seeker
// 	if err := db.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
// 		fmt.Println("Error fetching seeker data:", err) // Log error for debugging
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error fetching seeker data",
// 		})
// 		return
// 	}

// 	// Collect preferred titles
// 	var preferredTitles []string
// 	if seeker.PrimaryTitle != "" {
// 		preferredTitles = append(preferredTitles, seeker.PrimaryTitle)
// 	}
// 	// Check if there are preferred titles
// 	if len(preferredTitles) == 0 {
// 		fmt.Println("No preferred job titles set for user.") // Log for debugging
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "No preferred job titles set for user."})
// 		return
// 	}

// 	// Title filtering (case-insensitive)
// 	var conditions []string
// 	var values []interface{}
// 	for _, title := range preferredTitles {
// 		conditions = append(conditions, "LOWER(title) LIKE ?")
// 		values = append(values, "%"+strings.ToLower(title)+"%")
// 	}
// 	whereClause := strings.Join(conditions, " OR ")

// 	// Query jobs from the Seeker model (not saving, just fetching)
// 	var jobData []models.Job // Assuming `Job` model holds job-related data in Seeker model

// 	// Fetch jobs associated with the user (no saving here, just filtering)
// 	if err := db.Debug().Where(whereClause, values...).Find(&jobData).Error; err != nil {
// 		fmt.Println("Error fetching job data:", err) // Log error for debugging
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error fetching job data",
// 		})
// 		return
// 	}

// 	// Fetch user's professional summary to get skills (stored in Seeker model)
// 	var summary struct {
// 		Skills []string `json:"skills"`
// 	}
// 	if err := json.Unmarshal(seeker.ProfessionalSummary, &summary); err != nil {
// 		fmt.Println("Error unmarshalling professional summary:", err) // Log error for debugging
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to parse professional summary",
// 		})
// 		return
// 	}

// 	// Helper function to generate a random salary range (between 25k and 50k)
// 	randomSalary := func() dto.SalaryRange {
// 		minSalary := (rand.Intn(25) + 25) * 1000 // Random min salary between 25k and 50k
// 		maxSalary := (rand.Intn(25) + 25) * 1000 // Random max salary between 25k and 50k
// 		if minSalary > maxSalary {
// 			minSalary, maxSalary = maxSalary, minSalary // Ensure min is always less than max
// 		}
// 		return dto.SalaryRange{Min: minSalary, Max: maxSalary}
// 	}

// 	// Process job data stored in Seeker (no database saving)
// 	for _, job := range jobData {
// 		// Log job data for each job
// 		// fmt.Println("Processing job:", job)

// 		// Generate random salary range
// 		salaryRange := randomSalary()

// 		// Check the match score for this user and job pair
// 		var matchScore models.MatchScore
// 		if err := db.Where("seeker_id = ? AND job_id = ?", userID, job.JobID).First(&matchScore).Error; err != nil {
// 			// If no match score is found, set default to 50
// 			if err == gorm.ErrRecordNotFound {
// 				matchScore.MatchScore = 50
// 			} else {
// 				fmt.Println("Error fetching match score:", err) // Log error for debugging
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"error": "Error fetching match score",
// 				})
// 				return
// 			}
// 		}

// 		// Add job to jobs list with match score
// 		jobs = append(jobs, dto.JobDTO{
// 			Source:         "seeker", // All jobs are now stored under "seeker" in the Seeker model
// 			ID:             job.ID,   // Assuming job.ID is a UUID
// 			JobID:          job.JobID,
// 			Title:          job.Title,
// 			Company:        job.Company,
// 			Location:       job.Location,
// 			PostedDate:     job.PostedDate,
// 			Processed:      job.Processed,
// 			JobType:        job.JobType,
// 			Skills:         job.Skills, // Assuming Skills is part of the Job model stored in Seeker
// 			UserSkills:     summary.Skills, // Skills parsed from ProfessionalSummary
// 			ExpectedSalary: salaryRange,
// 			MatchScore:     matchScore.MatchScore, // Use the fetched or default match score
// 			Description:    job.JobDescription,    // Assuming Description is part of the Job model stored in Seeker
// 		})
// 	}

// 	// Respond with the filtered jobs and pagination
// 	c.JSON(http.StatusOK, gin.H{
// 		"jobs": jobs,
// 		"pagination": gin.H{
// 			"total": len(jobs),
// 		},
// 	})
// }
