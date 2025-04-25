package features

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"RAAS/models"
)
// GetNextEntryStep handles fetching the next incomplete step in the user entry timeline for MongoDB
func GetNextEntryStep() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get userID from context
		userID := c.MustGet("userID").(uuid.UUID)
		fmt.Println("UserID:", userID) // Debugging line

		// Get MongoDB database from context
		db := c.MustGet("db").(*mongo.Database)
		if db == nil {
			fmt.Println("Error: MongoDB database is nil")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database unavailable"})
			return
		}
		fmt.Println("MongoDB database successfully fetched") // Debugging line

		// Fetch the user entry timeline from the database
		collection := db.Collection("user_entry_timelines")
		var timeline models.UserEntryTimeline
		err := collection.FindOne(c, bson.M{"auth_user_id": userID}).Decode(&timeline)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("Error fetching timeline: User not found")
				c.JSON(http.StatusNotFound, gin.H{"error": "Timeline not found"})
			} else {
				fmt.Println("Error fetching timeline:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timeline"})
			}
			return
		}

		fmt.Println("Timeline data:", timeline)

		// Define the steps and their completion status
		steps := []struct {
			Name      string
			Completed bool
			Required  bool
		}{
			{"personal_infos", timeline.PersonalInfosCompleted, timeline.PersonalInfosRequired},
			{"professional_summaries", timeline.ProfessionalSummariesCompleted, timeline.ProfessionalSummariesRequired},
			{"work_experiences", timeline.WorkExperiencesCompleted, timeline.WorkExperiencesRequired},
			{"educations", timeline.EducationsCompleted, timeline.EducationsRequired},
			{"certificates", timeline.CertificatesCompleted, timeline.CertificatesRequired},
			{"languages", timeline.LanguagesCompleted, timeline.LanguagesRequired},
			{"preferred_job_titles", timeline.PreferredJobTitlesCompleted, timeline.PreferredJobTitlesRequired},
		}

		// Iterate through steps and check if any step is incomplete and required
		for _, step := range steps {
			fmt.Printf("Checking step: %s, Completed: %v, Required: %v\n", step.Name, step.Completed, step.Required)
			if step.Required && !step.Completed {
				c.JSON(http.StatusOK, gin.H{
					"completed": false,
					"next_step": step.Name,
				})
				return
			}
		}

		// If all required steps are complete, mark timeline as completed (if not already)
		if !timeline.Completed {
			fmt.Println("Marking timeline as completed")
			timeline.Completed = true
			update := bson.M{
				"$set": bson.M{"completed": true},
			}

			_, err := collection.UpdateOne(c, bson.M{"auth_user_id": userID}, update)
			if err != nil {
				fmt.Println("Error updating timeline:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as completed"})
				return
			}
			fmt.Println("Timeline marked as completed")
		}

		// Return the response
		c.JSON(http.StatusOK, gin.H{
			"completed": true,
			"next_step": nil,
		})
	}
}

