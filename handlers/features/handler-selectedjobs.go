package features

// import (
// 	"RAAS/models"
// 	"RAAS/dto"
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// 	"github.com/google/uuid"
// 	"net/http"
// 	"strings"
// )

// // SelectedJobsHandler struct for managing selected jobs
// type SelectedJobsHandler struct {
// 	DB *gorm.DB
// }

// // NewSelectedJobsHandler returns a new SelectedJobsHandler instance
// func NewSelectedJobsHandler(db *gorm.DB) *SelectedJobsHandler {
// 	return &SelectedJobsHandler{DB: db}
// }

// // GetSelectedJobs - Get all selected jobs for the authenticated user
// func (h *SelectedJobsHandler) GetSelectedJobs(c *gin.Context) {
// 	// Get user ID from JWT claims
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	var selectedJobs []models.SelectedJobApplication
// 	if err := h.DB.Where("auth_user_id = ?", userID).Find(&selectedJobs).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error fetching selected jobs",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"selected_jobs": selectedJobs,
// 	})
// }

// // PostSelectedJob - Create a new selected job for the authenticated user
// func (h *SelectedJobsHandler) PostSelectedJob(c *gin.Context) {
// 	// Get user ID from JWT claims
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	// Parse the request body to get the job details
// 	var jobDTO dto.JobDTO
// 	if err := c.ShouldBindJSON(&jobDTO); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid request payload",
// 		})
// 		return
// 	}

// 	// Fetch the Seeker record for the user
// 	var seeker models.Seeker
// 	if err := h.DB.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Seeker not found",
// 		})
// 		return
// 	}

// 	// Check if the Seeker has enough remaining daily selectable jobs
// 	if seeker.DailySelectableJobsCount <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "No available selectable jobs for today",
// 		})
// 		return
// 	}

// 	// Process the input for skills and userSkills (comma-separated lists)
// 	skillsList := strings.Split(jobDTO.Skills, ",")
// 	userSkillsList := strings.Join(jobDTO.UserSkills, ",") // assuming userSkills is a list of strings

// 	// Create the selected job record
// 	selectedJob := models.SelectedJobApplication{
// 		AuthUserID:            userID,
// 		JobID:                 jobDTO.JobID,
// 		Title:                 jobDTO.Title,
// 		Company:               jobDTO.Company,
// 		Location:              jobDTO.Location,
// 		PostedDate:           jobDTO.PostedDate,
// 		Processed:            jobDTO.Processed,
// 		JobType:              jobDTO.JobType,
// 		Skills:               strings.Join(skillsList, ", "), // Store skills as comma-separated string
// 		UserSkills:           userSkillsList,                 // Store user skills as comma-separated string
// 		MinSalary:            jobDTO.ExpectedSalary.Min,
// 		MaxSalary:            jobDTO.ExpectedSalary.Max,
// 		MatchScore:           jobDTO.MatchScore,
// 		Description:          jobDTO.Description,
// 		Source: 				jobDTO.Source,
		
// 	}

	
// 	// Insert the selected job into the database
// 	if err := h.DB.Create(&selectedJob).Error; err != nil {
// 		// Log the error for better debugging
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error creating selected job",
// 			"details": err.Error(),  // Log the specific database error
// 		})
// 		return
// 	}

// 	// Decrease the daily selectable jobs count for the Seeker
// 	seeker.DailySelectableJobsCount -= 1
// 	if err := h.DB.Save(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error updating seeker data",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":    "Job selected successfully",
// 		"selected_job": selectedJob,
// 	})
// }


// // UpdateSelectedJob - Update an existing selected job
// func (h *SelectedJobsHandler) UpdateSelectedJob(c *gin.Context) {
// 	// Get user ID from JWT claims
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	// Find the selected job by ID
// 	jobID := c.Param("id")
// 	var selectedJob models.SelectedJobApplication
// 	if err := h.DB.Where("auth_user_id = ? AND id = ?", userID, jobID).First(&selectedJob).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "Selected job not found",
// 		})
// 		return
// 	}

// 	// Parse the request body to get the updated job details
// 	var jobDTO dto.JobDTO
// 	if err := c.ShouldBindJSON(&jobDTO); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid request payload",
// 		})
// 		return
// 	}

// 	// Update the selected job record
// 	selectedJob.Title = jobDTO.Title
// 	selectedJob.Company = jobDTO.Company
// 	selectedJob.Location = jobDTO.Location
// 	selectedJob.Processed = jobDTO.Processed
// 	selectedJob.JobType = jobDTO.JobType
// 	selectedJob.Skills = jobDTO.Skills
// 	selectedJob.MatchScore = jobDTO.MatchScore
// 	selectedJob.Description = jobDTO.Description

// 	// Save the updated selected job
// 	if err := h.DB.Save(&selectedJob).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error updating selected job",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":    "Selected job updated successfully",
// 		"selected_job": selectedJob,
// 	})
// }

// // DeleteSelectedJob - Delete a selected job for the authenticated user
// func (h *SelectedJobsHandler) DeleteSelectedJob(c *gin.Context) {
// 	// Get user ID from JWT claims
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	// Find the selected job by ID
// 	jobID := c.Param("id")
// 	var selectedJob models.SelectedJobApplication
// 	if err := h.DB.Where("auth_user_id = ? AND id = ?", userID, jobID).First(&selectedJob).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "Selected job not found",
// 		})
// 		return
// 	}

// 	// Delete the selected job record
// 	if err := h.DB.Delete(&selectedJob).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error deleting selected job",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Selected job deleted successfully",
// 	})
// }
