package jobs

// import (
// 	"RAAS/internal/models"
// 	"net/http"
// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type JobDataHandler struct {

// }

// func NewJobDataHandler() *JobDataHandler {
// 	return &JobDataHandler{}
// }

// func (h *JobDataHandler) GetAllJobs(c *gin.Context) {
// 	var jobs []models.Job
// 	db := c.MustGet("db").(*mongo.Database)
// 	// Access the "jobs" collection in the MongoDB database
// 	jobCollection := db.Collection("jobs")

// 	// Fetch all jobs from the MongoDB collection
// 	cursor, err := jobCollection.Find(c, bson.M{})
// 	if err != nil {
// 		// Return an internal server error if fetching the jobs fails
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
// 		return
// 	}
// 	defer cursor.Close(c)

// 	// If no jobs are found
// 	if !cursor.Next(c) {
// 		c.JSON(http.StatusOK, gin.H{"message": "No jobs available"})
// 		return
// 	}

// 	// Iterate over the cursor and decode the documents into `jobs` slice
// 	for cursor.Next(c) {
// 		var job models.Job
// 		if err := cursor.Decode(&job); err != nil {
// 			// If error decoding, return internal server error
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode job data"})
// 			return
// 		}
// 		// Append the job to the jobs slice
// 		jobs = append(jobs, job)
// 	}

// 	// Return all jobs
// 	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
// }
