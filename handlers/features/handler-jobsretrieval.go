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

// JobRetrievalHandler handles the retrieval of jobs based on user's preferences and skills
func JobRetrievalHandler(c *gin.Context) {
	// Retrieve MongoDB database client and userID from the claims (context)
	db := c.MustGet("db").(*mongo.Database)  // MongoDB database from context
	userID := c.MustGet("userID").(string)    // User ID from claims in context (as string)

	// Define the MongoDB collections
	seekerCollection := db.Collection("seekers")
	jobCollection := db.Collection("jobs")

	var jobs []dto.JobDTO

	// Fetch user's Seeker data from the MongoDB collection
	var seeker models.Seeker
	if err := seekerCollection.FindOne(c, bson.M{"auth_user_id": userID}).Decode(&seeker); err != nil {
		fmt.Println("Error fetching seeker data:", err) // Log error for debugging
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

	// Extract skills from ProfessionalSummary (if present)
	var skills []string
	if seeker.ProfessionalSummary != nil {
		skills = extractSkills(seeker.ProfessionalSummary) // Use helper function to extract skills
	}

	// Title filtering (case-insensitive)
	var filter bson.M
	conditions := []bson.M{}
	for _, title := range preferredTitles {
		conditions = append(conditions, bson.M{"title": bson.M{"$regex": title, "$options": "i"}}) // Case-insensitive search
	}

	if len(conditions) > 0 {
		filter = bson.M{"$or": conditions}
	}

	// Use pagination middleware's `page` and `limit`
	page, limit := c.GetInt("page"), c.GetInt("limit")
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Calculate the skip value
	fmt.Println("Page and limit",page,limit)
	skip := (page - 1) * limit
	fmt.Println("Skip value:", skip)

	// Query jobs from MongoDB collection with pagination
	cursor, err := jobCollection.Find(c, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		fmt.Println("Error fetching job data:", err) // Log error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching job data",
		})
		return
	}
	defer cursor.Close(c)

	// Parse the job data into DTO
	for cursor.Next(c) {
		var job models.Job
		if err := cursor.Decode(&job); err != nil {
			fmt.Println("Error decoding job:", err) // Log error for debugging
			continue
		}

		// Generate a random salary range
		salaryRange := randomSalary()

		// Check the match score for this user and job pair
		var matchScore models.MatchScore
		if err := seekerCollection.FindOne(c, bson.M{"auth_user_id": userID, "job_id": job.JobID}).Decode(&matchScore); err != nil {
			// If no match score is found, set default to 50
			if err == mongo.ErrNoDocuments {
				matchScore.MatchScore = 50
			} else {
				fmt.Println("Error fetching match score:", err) // Log error for debugging
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error fetching match score",
				})
				return
			}
		}

		// Add job to jobs list with match score
		jobs = append(jobs, dto.JobDTO{
			Source:         "seeker", // All jobs are now stored under "seeker" in the Seeker model
			JobID:          job.JobID,
			Title:          job.Title,
			Company:        job.Company,
			Location:       job.Location,
			PostedDate:     job.PostedDate,
			Processed:      job.Processed,
			JobType:        job.JobType,
			Skills:         job.Skills, // Assuming Skills is part of the Job model stored in Seeker
			UserSkills:     skills, // Use extracted skills from the ProfessionalSummary
			ExpectedSalary: salaryRange,
			MatchScore:     matchScore.MatchScore, // Use the fetched or default match score
			Description:    job.JobDescription,    // Assuming Description is part of the Job model stored in Seeker
		})
	}

	// Get total count of jobs for pagination
	totalCount, err := jobCollection.CountDocuments(c, filter)
	if err != nil {
		fmt.Println("Error counting job data:", err) // Log error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error counting job data",
		})
		return
	}

		// Calculate next and previous page links
		nextPage := ""
		if int64(page)*int64(limit) < totalCount {
			nextPage = fmt.Sprintf("/api/jobs?page=%d&limit=%d", page+1, limit)
		}

	prevPage := ""
	if page > 1 {
		prevPage = fmt.Sprintf("/api/jobs?page=%d&limit=%d", page-1, limit)
	}

	// Respond with the filtered jobs, pagination, and links
	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
		"pagination": gin.H{
			"total":     totalCount,
			"next":      nextPage,
			"prev":      prevPage,
			"current":   page,
			"per_page":  limit,
		},
	})
}

// Helper function to generate a random salary range (between 25k and 50k)
func randomSalary() dto.SalaryRange {
	minSalary := (rand.Intn(25) + 25) * 1000 // Random min salary between 25k and 50k
	maxSalary := (rand.Intn(25) + 25) * 1000 // Random max salary between 25k and 50k
	if minSalary > maxSalary {
		minSalary, maxSalary = maxSalary, minSalary // Ensure min is always less than max
	}
	return dto.SalaryRange{Min: minSalary, Max: maxSalary}
}
