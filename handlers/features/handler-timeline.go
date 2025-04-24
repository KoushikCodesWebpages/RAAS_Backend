package features

import (
    "net/http"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "RAAS/models"
    "gorm.io/gorm"
)

func GetNextEntryStep(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get userID from context
        userID := c.MustGet("userID").(uuid.UUID)
        fmt.Println("UserID:", userID)  // Debugging line to check userID

        // Fetch the user entry timeline from the database
        var timeline models.UserEntryTimeline
        if err := db.First(&timeline, "auth_user_id = ?", userID).Error; err != nil {
            fmt.Println("Error fetching timeline:", err)  // Debugging line to log error
            c.JSON(http.StatusNotFound, gin.H{"error": "Timeline not found"})
            return
        }

        fmt.Println("Timeline data:", timeline)  // Debugging line to check the timeline data

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
            fmt.Println("Checking step:", step.Name, "Completed:", step.Completed, "Required:", step.Required)  // Debugging line

            if step.Required && !step.Completed {
                c.JSON(http.StatusOK, gin.H{
                    "completed": false,
                    "next_step": step.Name,
                })
                return
            }
        }

        // Mark completed if not already
        if !timeline.Completed {
            timeline.Completed = true
            if err := db.Save(&timeline).Error; err != nil {
                fmt.Println("Error saving updated timeline:", err)  // Debugging line to log error
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as completed"})
                return
            }
            fmt.Println("Timeline marked as completed")  // Debugging line to indicate timeline was marked completed
        }

        // Return the response
        c.JSON(http.StatusOK, gin.H{
            "completed": true,
            "next_step": nil,
        })
    }
}
