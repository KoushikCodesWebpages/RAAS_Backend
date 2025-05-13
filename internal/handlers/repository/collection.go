package repository

import (
	"context"
	"RAAS/internal/models"
	"RAAS/internal/dto"
	"RAAS/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// Fetch seeker and extract skills
func GetSeekerData(db *mongo.Database, userID string) (models.Seeker, []string, error) {
	var seeker models.Seeker
	err := db.Collection("seekers").FindOne(context.TODO(), bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		return models.Seeker{}, nil, err
	}
	skills := []string{}
	if seeker.ProfessionalSummary != nil {
		skills = ExtractSkills(seeker.ProfessionalSummary) // Use your existing skill extraction logic
	}
	return seeker, skills, nil
}


// Extract preferred titles from seeker
func CollectPreferredTitles(seeker models.Seeker) []string {
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

// Fetch job by job ID
func GetJobByID(db *mongo.Database, jobID string) (models.Job, error) {
	var job models.Job
	err := db.Collection("jobs").FindOne(context.TODO(), bson.M{"job_id": jobID}).Decode(&job)
	if err != nil {
		return models.Job{}, err
	}
	return job, nil
}

func CountJobsByTitles(db *mongo.Database, titles []string) (int64, error) {
	if len(titles) == 0 {
		return 0, fmt.Errorf("no titles provided for counting")
	}

	// Build the $or conditions for matching titles in primary, secondary, tertiary title fields (case-insensitive)
	var orConditions []bson.M
	for _, title := range titles {
		regexFilter := bson.M{
			"$regex":   title,
			"$options": "i", // Case-insensitive match
		}
		orConditions = append(orConditions, bson.M{"primary_title": regexFilter})
		orConditions = append(orConditions, bson.M{"secondary_title": regexFilter})
		orConditions = append(orConditions, bson.M{"tertiary_title": regexFilter})
	}

	// Final filter: $or condition
	filter := bson.M{
		"$or": orConditions,
	}

	log.Printf("[DEBUG] Counting jobs with filter: %+v", filter)

	// Perform the count
	count, err := db.Collection("jobs").CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Printf("[ERROR] Failed to count jobs by titles: %v", err)
		return 0, err
	}

	log.Printf("[DEBUG] Found %d jobs matching given titles", count)
	return count, nil
}



// Extract skills safely
func ExtractSkills(professionalSummary bson.M) []string {
	if val, ok := professionalSummary["skills"].(primitive.A); ok {
		var skills []string
		for _, skill := range val {
			if str, ok := skill.(string); ok {
				skills = append(skills, str)
			}
		}
		return skills
	}
	return nil
}

// Helper function to fetch saved job IDs
func FetchSavedJobIDs(c *gin.Context, col *mongo.Collection, userID string) ([]string, error) {
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

// ConvertBsonMToWorkExperience converts bson.M to WorkExperienceResponse
func ConvertBsonMToWorkExperience(data bson.M) (dto.WorkExperienceResponse, error) {
	var workExp dto.WorkExperienceResponse

	// Convert start_date from time.Time to utils.DateOnly
	startDate, ok := data["start_date"].(time.Time)
	if !ok {
		log.Printf("[ERROR] Invalid or missing start_date: %+v", data["start_date"])
		return workExp, fmt.Errorf("invalid or missing start_date")
	}

	// Convert startDate (time.Time) to DateOnly (utils.DateOnly)
	workExp.StartDate = utils.ToDateOnly(startDate)

	// Handle optional end_date (convert if present)
	if endDate, ok := data["end_date"].(time.Time); ok {
		endDateOnly := utils.ToDateOnly(endDate)
		workExp.EndDate = &endDateOnly
	}

	return workExp, nil
}


func GetExperienceInMonths(workExperiences []bson.M) (int, error) {
	var totalMonths int

	for i, exp := range workExperiences {
		// Convert bson.M to WorkExperienceResponse
		workExp, err := ConvertBsonMToWorkExperience(exp)
		if err != nil {
			return 0, fmt.Errorf("error converting bson.M to WorkExperience: %w", err)
		}

		// Log the entire work experience entry
		log.Printf("[DEBUG] Work experience #%d: %+v", i+1, workExp)

		// Convert workExp.StartDate (utils.DateOnly) to time.Time for calculation
		startTime := utils.ToTime(workExp.StartDate)
		if startTime.IsZero() {
			log.Printf("[WARN] Missing or invalid start_date for experience #%d", i+1)
			// Optionally, return error or skip this entry based on your business logic
			continue
		}

		// Default end time to current time if not provided
		endTime := time.Now()

		// Extract and validate the end date, if present
		if workExp.EndDate != nil {
			endTime = utils.ToTime(*workExp.EndDate)
			if endTime.IsZero() {
				log.Printf("[WARN] Invalid end_date for experience #%d; defaulting to current time", i+1)
			}
		}

		// Calculate the duration in months
		years := endTime.Year() - startTime.Year()
		months := int(endTime.Month()) - int(startTime.Month())
		durationInMonths := years*12 + months

		// Ensure non-negative duration
		if durationInMonths < 0 {
			log.Printf("[WARN] Negative duration calculated for experience #%d; setting to 0", i+1)
			durationInMonths = 0
		}

		log.Printf("[DEBUG] Experience #%d duration: %d months", i+1, durationInMonths)

		totalMonths += durationInMonths
	}

	log.Printf("[DEBUG] Total experience across all entries: %d months", totalMonths)
	return totalMonths, nil
}



// Helper function to extract certificates
func ExtractCertificates(certificates []bson.M) []string {
	var result []string
	for _, cert := range certificates {
		if certName, ok := cert["certificate_name"].(string); ok {
			result = append(result, certName)
		}
	}
	return result
}

