package features

import (
	"RAAS/dto"
	"RAAS/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SelectedJobsHandler handles job selection operations
type SelectedJobsHandler struct{}

// NewSelectedJobsHandler initializes and returns a new instance of SelectedJobsHandler
func NewSelectedJobsHandler() *SelectedJobsHandler {
	return &SelectedJobsHandler{}
}

// Random salary range generator
func randomSalary() (int, int) {
	// Example random salary range logic, adjust as needed
	minSalary := 50000 // Example minimum salary
	maxSalary := 100000 // Example maximum salary
	return minSalary, maxSalary
}
// PostSelectedJob saves the selected job for the authenticated user
func (h *SelectedJobsHandler) PostSelectedJob(c *gin.Context) {
	// Get the database from the context
	db := c.MustGet("db").(*mongo.Database)
	jobCollection := db.Collection("jobs")   // Collection where job data is stored
	selectedJobsCollection := db.Collection("selected_job_applications") // Collection where selected jobs are stored
	seekerCollection := db.Collection("seekers") // Collection where seeker data is stored

	// Retrieve the user ID from the context
	userID := c.MustGet("userID").(string)

	// Input: Expect a JobID as input to retrieve the job data
	var input struct {
		JobID string `json:"job_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	// Check if the job has already been selected by the user
	var existingSelection models.SelectedJobApplication
	err := selectedJobsCollection.FindOne(c, bson.M{"auth_user_id": userID, "job_id": input.JobID}).Decode(&existingSelection)
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusConflict, gin.H{
			"error": "You have already selected this job",
		})
		return
	}

	// Retrieve the job data based on the jobID
	var job models.Job
	if err := jobCollection.FindOne(c, bson.M{"job_id": input.JobID}).Decode(&job); err != nil {
		fmt.Println("Error retrieving job data:", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
		return
	}

	// Retrieve the seeker data based on the authenticated user
	var seeker models.Seeker
	if err := seekerCollection.FindOne(c, bson.M{"auth_user_id": userID}).Decode(&seeker); err != nil {
		fmt.Println("Error fetching seeker data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching seeker data",
		})
		return
	}

	// Map the professional summary to skills
	var skills []string
	if seeker.ProfessionalSummary != nil {
		skills = extractSkills(seeker.ProfessionalSummary)
	}

	// Generate a random salary range
	minSalary, maxSalary := randomSalary()

	// Create SalaryRange struct
	expectedSalary := models.SalaryRange{
		Min: minSalary,
		Max: maxSalary,
	}

	// Prepare the selected job data to save into the selected_jobs collection
	selectedJob := models.SelectedJobApplication{
		AuthUserID:            userID,
		JobID:                 job.JobID,
		Title:                 job.Title,
		Company:               job.Company,
		Location:              job.Location,
		PostedDate:            job.PostedDate,
		Processed:             true, // Set processed to always true
		JobType:               job.JobType,
		Skills:                job.Skills,
		UserSkills:            skills, // Populating user skills from seeker profile
		ExpectedSalary:        expectedSalary, // Set SalaryRange
		MatchScore:            80, // Always set to 80
		Description:           job.JobDescription,
		Selected:              true,
		CvGenerated:           false, // Or set to true if the CV was generated
		CoverLetterGenerated:  false, // Or set to true if cover letter was generated
		ViewLink:              false, // Set view link true based on your business logic
		SelectedDate:          time.Now(),
	}

	// Insert the selected job data into the selected_jobs collection
	_, err = selectedJobsCollection.InsertOne(c, selectedJob)
	if err != nil {
		fmt.Println("Error saving selected job:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save selected job",
		})
		return
	}

	// Update the Seeker's DailySelectableJobsCount
	updateSeeker := bson.M{
		"$inc": bson.M{
			"daily_selectable_jobs_count": -1, // Reduce the count by 1
		},
	}

	_, err = seekerCollection.UpdateOne(c, bson.M{"auth_user_id": userID}, updateSeeker)
	if err != nil {
		fmt.Println("Error updating seeker data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update seeker data",
		})
		return
	}

	// Increment the SelectedCount in the Job model
	updateJob := bson.M{
		"$inc": bson.M{
			"selected_count": 1, // Increment the selected count by 1
		},
	}

	_, err = jobCollection.UpdateOne(c, bson.M{"job_id": job.JobID}, updateJob)
	if err != nil {
		fmt.Println("Error updating job data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update job data",
		})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{
		"message": "Selected job saved successfully",
	})
}

func (h *SelectedJobsHandler) GetSelectedJobs(c *gin.Context) {
	// Get the database from the context
	db := c.MustGet("db").(*mongo.Database)
	selectedJobsCollection := db.Collection("selected_job_applications") // Collection where selected jobs are stored

	// Retrieve the user ID from the context
	userID := c.MustGet("userID").(string)

	// Define the filter to fetch selected jobs for the authenticated user
	filter := bson.M{"auth_user_id": userID}

	// Access pagination values from context set by middleware
	pagination := c.MustGet("pagination").(gin.H)
	offsetInt := pagination["offset"].(int)
	limitInt := pagination["limit"].(int)

	// Define the pagination options
	findOptions := options.Find().SetSkip(int64(offsetInt)).SetLimit(int64(limitInt))

	// Query the database to retrieve selected jobs for the authenticated user
	cursor, err := selectedJobsCollection.Find(c, filter, findOptions)
	if err != nil {
		fmt.Println("Error fetching selected jobs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching selected jobs",
		})
		return
	}
	defer cursor.Close(c)

	// Slice to hold the response data
	var selectedJobs []dto.SelectedJobResponse

	// Iterate through the cursor and decode the documents
	for cursor.Next(c) {
		var selectedJob models.SelectedJobApplication
		if err := cursor.Decode(&selectedJob); err != nil {
			fmt.Println("Error decoding selected job:", err)
			continue
		}

		// Convert SelectedJobApplication to SelectedJobResponse DTO
		selectedJobResponse := dto.SelectedJobResponse{
			AuthUserID:            selectedJob.AuthUserID,
			Source:                selectedJob.Source,
			JobID:                 selectedJob.JobID,
			Title:                 selectedJob.Title,
			Company:               selectedJob.Company,
			Location:              selectedJob.Location,
			PostedDate:            selectedJob.PostedDate,
			Processed:             selectedJob.Processed,
			JobType:               selectedJob.JobType,
			Skills:                selectedJob.Skills,
			UserSkills:            selectedJob.UserSkills,
			ExpectedSalary:        convertSalaryRange(selectedJob.ExpectedSalary),
			MatchScore:            selectedJob.MatchScore,
			Description:           selectedJob.Description,
			Selected:              selectedJob.Selected,
			CvGenerated:           selectedJob.CvGenerated,
			CoverLetterGenerated:  selectedJob.CoverLetterGenerated,
			ViewLink:              selectedJob.ViewLink,
			SelectedDate:          selectedJob.SelectedDate.Format(time.RFC3339), // Formatting SelectedDate to string
		}

		// Append the DTO to the response slice
		selectedJobs = append(selectedJobs, selectedJobResponse)
	}

	// If no selected jobs were found
	if len(selectedJobs) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "No selected jobs found",
		})
		return
	}

	// Count total documents for pagination (to calculate total pages)
	totalCount, err := selectedJobsCollection.CountDocuments(c, filter)
	if err != nil {
		fmt.Println("Error counting selected jobs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error counting selected jobs",
		})
		return
	}

	// Create next and prev page URLs for pagination
	nextPage := ""
	if int64(offsetInt+limitInt) < totalCount {
		nextPage = fmt.Sprintf("/api/selected-jobs?offset=%d&limit=%d", offsetInt+limitInt, limitInt)
	}

	prevPage := ""
	if offsetInt > 0 {
		prevPage = fmt.Sprintf("/api/selected-jobs?offset=%d&limit=%d", offsetInt-limitInt, limitInt)
	}

	// Send JSON response with selected jobs and pagination info
	c.JSON(http.StatusOK, gin.H{
		"selected_jobs": selectedJobs,
		"pagination": gin.H{
			"total":     totalCount,
			"next":      nextPage,
			"prev":      prevPage,
			"current":   (offsetInt / limitInt) + 1,
			"per_page":  limitInt,
		},
	})
}

func convertSalaryRange(modelSalary models.SalaryRange) dto.SalaryRange {
	return dto.SalaryRange{
		Min: modelSalary.Min,
		Max: modelSalary.Max,
	}
}