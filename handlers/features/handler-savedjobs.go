package features
import (
	"RAAS/dto"
	"RAAS/models"
	"RAAS/handlers"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SavedJobsHandler handles saving and retrieving saved jobs
type SavedJobsHandler struct{}

// NewSavedJobsHandler initializes and returns a new SavedJobsHandler instance
func NewSavedJobsHandler() *SavedJobsHandler {
	return &SavedJobsHandler{}
}

// POST /saved-jobs
func (h *SavedJobsHandler) SaveJob(c *gin.Context) {
	db := c.MustGet("db").(*mongo.Database)
	userID := c.MustGet("userID").(string)

	var payload struct {
		JobID  string `json:"job_id" binding:"required"`
		Source string `json:"source"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	savedJob := models.SavedJob{
		AuthUserID: userID,
		Source:     payload.Source,
		JobID:      payload.JobID,
	}

	_, err := db.Collection("saved_jobs").InsertOne(c, savedJob)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Job already saved"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save job"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Job saved successfully"})
}

// GET /saved-jobs
func (h *SavedJobsHandler) GetSavedJobs(c *gin.Context) {
	db := c.MustGet("db").(*mongo.Database)
	userID := c.MustGet("userID").(string)

	savedJobIDs, err := fetchSavedJobIDs(c, db.Collection("saved_jobs"), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching saved jobs"})
		return
	}

	if len(savedJobIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"pagination": gin.H{
				"total":    0,
				"next":     "",
				"prev":     "",
				"current":  1,
				"per_page": 10,
			},
			"jobs": []dto.JobDTO{},
		})
		return
	}

	pagination := c.MustGet("pagination").(gin.H)
	offset := pagination["offset"].(int)
	limit := pagination["limit"].(int)

	filter := bson.M{"job_id": bson.M{"$in": savedJobIDs}}
	cursor, err := db.Collection("jobs").Find(c, filter, options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching job data"})
		return
	}
	defer cursor.Close(c)

	_, skills, err := handlers.GetSeekerData(db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user data"})
		return
	}

	var jobs []dto.JobDTO
	for cursor.Next(c) {
		var job models.Job
		if err := cursor.Decode(&job); err != nil {
			continue
		}

		expectedSalary := dto.SalaryRange(handlers.GenerateSalaryRange())
		jobs = append(jobs, dto.JobDTO{
			Source:         "saved",
			JobID:          job.JobID,
			Title:          job.Title,
			Company:        job.Company,
			Location:       job.Location,
			PostedDate:     job.PostedDate,
			Processed:      job.Processed,
			JobType:        job.JobType,
			Skills:         job.Skills,
			UserSkills:     skills,
			ExpectedSalary: expectedSalary,
			MatchScore:     50,
			Description:    job.JobDescription,
		})
	}

	totalCount := int64(len(savedJobIDs))
	nextPage := ""
	if int64(offset+limit) < totalCount {
		nextPage = fmt.Sprintf("/api/saved-jobs?offset=%d&limit=%d", offset+limit, limit)
	}
	prevPage := ""
	if offset > 0 {
		prevPage = fmt.Sprintf("/api/saved-jobs?offset=%d&limit=%d", offset-limit, limit)
	}

	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":    totalCount,
			"next":     nextPage,
			"prev":     prevPage,
			"current":  (offset / limit) + 1,
			"per_page": limit,
		},
		"jobs": jobs,
	})
}

// Helper function to fetch saved job IDs
func fetchSavedJobIDs(c *gin.Context, col *mongo.Collection, userID string) ([]string, error) {
	var jobIDs []string
	cursor, err := col.Find(c, bson.M{"auth_user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var saved models.SavedJob
		if err := cursor.Decode(&saved); err == nil {
			jobIDs = append(jobIDs, saved.JobID)
		}
	}
	return jobIDs, nil
}
