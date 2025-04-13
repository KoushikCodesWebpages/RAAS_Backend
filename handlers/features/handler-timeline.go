package features

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "RAAS/models"
    "gorm.io/gorm"
)
func GetNextEntryStep(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("userID").(uuid.UUID)

        var timeline models.UserEntryTimeline
        if err := db.First(&timeline, "user_id = ?", userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Timeline not found"})
            return
        }

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

        for _, step := range steps {
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
            db.Save(&timeline)
        }

        c.JSON(http.StatusOK, gin.H{
            "completed": true,
            "next_step": nil,
        })
    }
}