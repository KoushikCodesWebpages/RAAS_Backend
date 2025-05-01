package features

import (
	"RAAS/dto"
	"RAAS/handlers"
	"RAAS/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func JobRetrievalHandler(c *gin.Context) {
	db := c.MustGet("db").(*mongo.Database)
	userID := c.MustGet("userID").(string)

	// --- Fetch seeker and skills using helper ---
	seeker, skills, err := handlers.GetSeekerData(db, userID)
	if err != nil {
		fmt.Println("Error fetching seeker data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching seeker data"})
		return
	}

	if seeker.PrimaryTitle == "" {
		c.JSON(http.StatusNoContent, gin.H{"error": "No preferred job title set for user."})
		return
	}

	// --- Collect preferred titles using helper ---
	preferredTitles := collectPreferredTitles(seeker)
	if len(preferredTitles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No preferred job titles set for user."})
		return
	}

	// --- Get applied jobs using helper ---
	appliedJobIDs, err := fetchAppliedJobIDs(c, db.Collection("selected_job_applications"), userID)
	if err != nil {
		fmt.Println("Error fetching applied jobs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching applied job data"})
		return
	}

	// --- Ensure appliedJobIDs is not nil ---
	if appliedJobIDs == nil {
		appliedJobIDs = []string{}
	}

	// --- Build MongoDB query ---
	filter := buildJobFilter(preferredTitles, appliedJobIDs)

	// --- Pagination ---
	pagination := c.MustGet("pagination").(gin.H)
	offset := pagination["offset"].(int)
	limit := pagination["limit"].(int)

	// --- Query jobs ---
	cursor, err := db.Collection("jobs").Find(c, filter, options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)))
	if err != nil {
		fmt.Println("Error fetching job data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching job data"})
		return
	}
	defer cursor.Close(c)

	var jobs []dto.JobDTO
	for cursor.Next(c) {
		var job models.Job
		if err := cursor.Decode(&job); err != nil {
			fmt.Println("Error decoding job:", err)
			continue
		}

		expectedSalary := dto.SalaryRange(handlers.GenerateSalaryRange())

		jobs = append(jobs, dto.JobDTO{
			Source:         "seeker",
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

	// --- Count total jobs ---
	totalCount, err := db.Collection("jobs").CountDocuments(c, filter)
	if err != nil {
		fmt.Println("Error counting job data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting job data"})
		return
	}

	// --- Build pagination response ---
	nextPage := ""
	if int64(offset+limit) < totalCount {
		nextPage = fmt.Sprintf("/api/jobs?offset=%d&limit=%d", offset+limit, limit)
	}
	prevPage := ""
	if offset > 0 {
		prevPage = fmt.Sprintf("/api/jobs?offset=%d&limit=%d", offset-limit, limit)
	}

	// --- Send response ---
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

// Extract preferred titles from seeker
func collectPreferredTitles(seeker models.Seeker) []string {
	var titles []string
	if seeker.PrimaryTitle != "" {
		titles = append(titles, seeker.PrimaryTitle)
	}
	if seeker.SecondaryTitle != nil && *seeker.SecondaryTitle != "" {
		titles = append(titles, *seeker.SecondaryTitle)
	}
	if seeker.TertiaryTitle != nil && *seeker.TertiaryTitle != "" {
		titles = append(titles, *seeker.TertiaryTitle)
	}
	return titles
}

// Fetch applied job IDs to exclude
func fetchAppliedJobIDs(c *gin.Context, col *mongo.Collection, userID string) ([]string, error) {
	jobIDs := []string{} // Always return a non-nil slice

	cursor, err := col.Find(c, bson.M{"auth_user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var application models.SelectedJobApplication
		if err := cursor.Decode(&application); err == nil {
			jobIDs = append(jobIDs, application.JobID)
		}
	}
	return jobIDs, nil
}

// Construct the job query filter
func buildJobFilter(preferredTitles, appliedJobIDs []string) bson.M {
	var titleConditions []bson.M
	for _, title := range preferredTitles {
		titleConditions = append(titleConditions, bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	}

	filter := bson.M{
		"$and": []bson.M{
			{"$or": titleConditions},
			{"job_id": bson.M{"$nin": appliedJobIDs}}, // safe now
		},
	}
	return filter
}
