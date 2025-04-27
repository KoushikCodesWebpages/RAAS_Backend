package features

import (
	"RAAS/dto"
	"RAAS/models"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
)


func JobRetrievalHandler(c *gin.Context) {

	db := c.MustGet("db").(*mongo.Database)  
	userID := c.MustGet("userID").(string)

	// Define the MongoDB collections
	seekerCollection := db.Collection("seekers")
	jobCollection := db.Collection("jobs")

	var jobs []dto.JobDTO

	// Fetch user's Seeker data from the MongoDB collection
	var seeker models.Seeker
	if err := seekerCollection.FindOne(c, bson.M{"auth_user_id": userID}).Decode(&seeker); err != nil {
		fmt.Println("Error fetching seeker data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching seeker data",
		})
		return
	}

	// Check if PrimaryTitle is empty or nil
	if seeker.PrimaryTitle == "" {
		// If PrimaryTitle is empty, return 204 No Content
		c.JSON(http.StatusNoContent, gin.H{"error": "No preferred job title set for user."})
		return
	}

	// Collect preferred titles from the Seeker model
	var preferredTitles []string
	if seeker.PrimaryTitle != "" {
		preferredTitles = append(preferredTitles, seeker.PrimaryTitle)
	}
	if seeker.SecondaryTitle != nil && *seeker.SecondaryTitle != "" {
		preferredTitles = append(preferredTitles, *seeker.SecondaryTitle)
	}
	if seeker.TertiaryTitle != nil && *seeker.TertiaryTitle != "" {
		preferredTitles = append(preferredTitles, *seeker.TertiaryTitle)
	}

	// Check if there are preferred titles
	if len(preferredTitles) == 0 {
		fmt.Println("No preferred job titles set for user.") // Log for debugging
		c.JSON(http.StatusBadRequest, gin.H{"error": "No preferred job titles set for user."})
		return
	}
	var skills []string
	if seeker.ProfessionalSummary != nil {
		skills = extractSkills(seeker.ProfessionalSummary)
	}
	var filter bson.M
	conditions := []bson.M{}
	for _, title := range preferredTitles {
		conditions = append(conditions, bson.M{"title": bson.M{"$regex": title, "$options": "i"}}) // Case-insensitive search
	}

	if len(conditions) > 0 {
		filter = bson.M{"$or": conditions}
	}

	pagination := c.MustGet("pagination").(gin.H)
	offset := pagination["offset"].(int)
	limit := pagination["limit"].(int)

	cursor, err := jobCollection.Find(c, filter, options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)))
	if err != nil {
		fmt.Println("Error fetching job data:", err) 
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching job data",
		})
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var job models.Job
		if err := cursor.Decode(&job); err != nil {
			fmt.Println("Error decoding job:", err) 
			continue
		}
		salaryRange := randomSalary()
		var matchScore models.MatchScore
		if err := seekerCollection.FindOne(c, bson.M{"auth_user_id": userID, "job_id": job.JobID}).Decode(&matchScore); err != nil {
			if err == mongo.ErrNoDocuments {
				matchScore.MatchScore = 50
			} else {
				fmt.Println("Error fetching match score:", err) 
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error fetching match score",
				})
				return
			}
		}

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
			ExpectedSalary: salaryRange,
			MatchScore:     matchScore.MatchScore, 
			Description:    job.JobDescription,   
		})
	}

	// Get total count of jobs for pagination
	totalCount, err := jobCollection.CountDocuments(c, filter)
	if err != nil {
		fmt.Println("Error counting job data:", err) 
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error counting job data",
		})
		return
	}

	nextPage := ""
	if int64(offset+limit) < totalCount {
		nextPage = fmt.Sprintf("/api/jobs?offset=%d&limit=%d", offset+limit, limit)
	}

	prevPage := ""
	if offset > 0 {
		prevPage = fmt.Sprintf("/api/jobs?offset=%d&limit=%d", offset-limit, limit)
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
		"pagination": gin.H{
			"total":     totalCount,
			"next":      nextPage,
			"prev":      prevPage,
			"current":   (offset / limit) + 1,
			"per_page":  limit,
		},
	})
}

func randomSalary() dto.SalaryRange {
	minSalary := (rand.Intn(25) + 25) * 1000 
	maxSalary := (rand.Intn(25) + 25) * 1000 
	if minSalary > maxSalary {
		minSalary, maxSalary = maxSalary, minSalary 
	}
	return dto.SalaryRange{Min: minSalary, Max: maxSalary}
}
