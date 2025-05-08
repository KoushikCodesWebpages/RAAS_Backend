package features

import (

	"RAAS/internal/dto"
	"RAAS/internal/models"

	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"

)

type SeekerProfileHandler struct{}

func NewSeekerProfileHandler() *SeekerProfileHandler {
	return &SeekerProfileHandler{}
}

// GetSeekerProfile retrieves the profile for the authenticated user
func (h *SeekerProfileHandler) GetSeekerProfile(c *gin.Context) {
	// Get authenticated user ID and db from context
	userID := c.MustGet("userID").(string)
	db := c.MustGet("db").(*mongo.Database)

	seekersCollection := db.Collection("seekers")

	// Set a timeout for the MongoDB operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the seeker by auth_user_id
	var seeker models.Seeker
	err := seekersCollection.FindOne(ctx, bson.M{"auth_user_id": userID}).Decode(&seeker)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seeker profile not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving seeker profile"})
		}
		log.Printf("Error retrieving seeker profile for auth_user_id: %s, Error: %v", userID, err)
		return
	}

	var languageNames []string
	for _, language := range seeker.Languages {
		// Ensure 'language' is a map (bson.M), and access the "language" key
		if lang, ok := language["language"].(string); ok {
			languageNames = append(languageNames, lang) // Collect the language name
		} else {
			log.Printf("[WARN] Invalid or missing 'language' field in languages array")
		}
	}
	// Map seeker to Seeker	ProfileDTO
	profile := dto.SeekerProfileDTO{
		AuthUserID:                  seeker.AuthUserID,
		FirstName:                   dereferenceString(getOptionalField(seeker.PersonalInfo, "first_name")),
		SecondName:                  getOptionalField(seeker.PersonalInfo, "second_name"),
		Skills:                      extractSkills(seeker.ProfessionalSummary),
		TotalExperienceInMonths:     getExperienceInMonths(seeker.WorkExperiences),
		Certificates:                extractCertificates(seeker.Certificates),
		PreferredJobTitle:           seeker.PrimaryTitle,
		SubscriptionTier:            seeker.SubscriptionTier,
		DailySelectableJobsCount:    seeker.DailySelectableJobsCount,
		DailyGeneratableCV:          seeker.DailyGeneratableCV,
		DailyGeneratableCoverletter: seeker.DailyGeneratableCoverletter,
		TotalApplications:           seeker.TotalApplications,
		TotalJobsAvailable:          0, // For now, as you said
		ProfileCompletion:           calculateProfileCompletion(seeker),
		Languages:                  languageNames, 
	}

	c.JSON(http.StatusOK, profile)
}



func dereferenceString(str *string) string {
	if str != nil {
		return *str
	}
	return "" // Return an empty string if the pointer is nil
}


// Helper function to get optional fields
func getOptionalField(info bson.M, field string) *string {
	if val, ok := info[field]; ok && val != nil {
		v := val.(string)
		return &v
	}
	return nil
}

// Extract skills safely
func extractSkills(professionalSummary bson.M) []string {
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

func getExperienceInMonths(workExperiences []bson.M) int {
	totalMonths := 0
	for i, exp := range workExperiences {
		// Print the entire work experience first
		log.Printf("[DEBUG] Work experience #%d: %v", i+1, exp)

		if startDate, ok := exp["start_date"].(map[string]interface{}); ok {
			// Directly extract the Unix timestamp (in milliseconds) from the start_date map
			startTimeUnixMillis, startOk := startDate["time"].(float64)
			if !startOk {
				log.Printf("[WARN] No start time found for experience #%d", i+1)
				continue // Skip if there's no start time
			}

			// Convert the Unix timestamp (in milliseconds) to seconds
			startTimeUnixSecs := time.Unix(int64(startTimeUnixMillis/1000), 0)

			// Log the parsed start time
			log.Printf("[DEBUG] Start time for experience #%d: %s", i+1, startTimeUnixSecs)

			// Default to current time if there's no end date
			end := time.Now()

			// Check if end date exists, if it does, use it
			if endDate, ok := exp["end_date"].(map[string]interface{}); ok && endDate["time"] != nil {
				endTimeUnixMillis, endOk := endDate["time"].(float64)
				if endOk {
					// Convert the Unix timestamp (in milliseconds) to seconds
					end = time.Unix(int64(endTimeUnixMillis/1000), 0)
				}
			}

			// Log the parsed end time
			log.Printf("[DEBUG] End time for experience #%d: %s", i+1, end)

			// Calculate duration in months
			years := end.Year() - startTimeUnixSecs.Year()
			months := int(end.Month()) - int(startTimeUnixSecs.Month())

			durationInMonths := years*12 + months
			if durationInMonths < 0 {
				durationInMonths = 0 // just in case
			}

			// Debugging log
			log.Printf("[DEBUG] Experience #%d duration: %d months", i+1, durationInMonths)

			// Add to total
			totalMonths += durationInMonths
		} else {
			// Log if start date is missing
			log.Printf("[WARN] No start date found for experience #%d", i+1)
		}
	}

	// Debugging total experience
	log.Printf("[DEBUG] Total experience across all entries: %d months", totalMonths)
	return totalMonths
}




// Helper function to extract certificates
func extractCertificates(certificates []bson.M) []string {
	var result []string
	for _, cert := range certificates {
		if certName, ok := cert["certificate_name"].(string); ok {
			result = append(result, certName)
		}
	}
	return result
}

// Helper function to calculate profile completion
func calculateProfileCompletion(seeker models.Seeker) int {
	completion := 0

	// Personal Info
	if seeker.PersonalInfo != nil {
		if seeker.PersonalInfo["first_name"] != nil {
			completion += 10
		}
		if seeker.PersonalInfo["second_name"] != nil {
			completion += 10
		}
	}

	// Skills
	if skills := extractSkills(seeker.ProfessionalSummary); len(skills) > 0 {
		completion += 20
	}

	// Work Experience
	if len(seeker.WorkExperiences) > 0 {
		completion += 20
	}

	// Certificates
	if len(seeker.Certificates) > 0 {
		completion += 20
	}

	// Preferred Job Title
	if seeker.PrimaryTitle != "" {
		completion += 20
	}

	// Subscription Tier
	if seeker.SubscriptionTier != "" {
		completion += 10
	}

	// Ensure completion is capped at 100
	if completion > 100 {
		completion = 100
	}

	return completion
}
