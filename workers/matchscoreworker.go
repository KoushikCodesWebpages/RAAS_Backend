package workers

import (
	"fmt"
	"log"
	"RAAS/models"
	"gorm.io/gorm"
	"time"
	"strings"
	"github.com/google/uuid"
)

// MatchScoreWorker continuously calculates match scores and stores them
type MatchScoreWorker struct {
	DB *gorm.DB
}

func (w *MatchScoreWorker) calculateAndStoreMatchScore(seekerAuthUserID uuid.UUID, jobID string, jobType string) error {
	// Fetch seeker details from the DB
	var seeker models.Seeker
	//log.Printf("Fetching seeker details for AuthUserID: %s", seekerAuthUserID)
	if err := w.DB.Where("auth_user_id = ?", seekerAuthUserID).First(&seeker).Error; err != nil {
		return fmt.Errorf("failed to find seeker: %v", err)
	}
	//log.Printf("Seeker details fetched for %s: %+v", seekerAuthUserID, seeker)

	var job interface{}
	if jobType == "linkedin" {
		var linkedinJob models.LinkedInJobMetaData
		//log.Printf("Fetching LinkedIn job details for JobID: %s", jobID)
		if err := w.DB.Where("job_id = ?", jobID).First(&linkedinJob).Error; err != nil {
			return fmt.Errorf("failed to find LinkedIn job: %v", err)
		}
		job = linkedinJob
		//log.Printf("LinkedIn job details fetched: %+v", linkedinJob)
	} else if jobType == "xing" {
		var xingJob models.XingJobMetaData
		//log.Printf("Fetching Xing job details for JobID: %s", jobID)
		if err := w.DB.Where("job_id = ?", jobID).First(&xingJob).Error; err != nil {
			return fmt.Errorf("failed to find Xing job: %v", err)
		}
		job = xingJob
		//log.Printf("Xing job details fetched: %+v", xingJob)
	}

	// Calculate the match score using the static function
	//log.Printf("Calculating match score for SeekerID: %s and JobID: %s", seekerAuthUserID, jobID)
	matchScore, err := CalculateMatchScore(seeker, job)
	if err != nil {
		return fmt.Errorf("failed to calculate match score: %v", err)
	}
	log.Printf("Calculated match score: %f", matchScore)

	// Create the MatchScore entry
	matchScoreEntry := models.MatchScore{
		SeekerID:   seeker.AuthUserID, // Use AuthUserID (uuid.UUID) for SeekerID
		JobID:      jobID,             // Use JobID from LinkedIn or Xing Job
		MatchScore: matchScore,
	}

	// Save the match score to the DB
	//log.Printf("Saving match score to DB for SeekerID: %s and JobID: %s", seekerAuthUserID, jobID)
	if err := w.DB.Create(&matchScoreEntry).Error; err != nil {
		return fmt.Errorf("failed to save match score: %v", err)
	}
	//log.Printf("Match score successfully saved for SeekerID: %s and JobID: %s", seekerAuthUserID, jobID)

	return nil
}

func (w *MatchScoreWorker) Run() {
	for {
		//log.Println("MatchScoreWorker started...")

		// Fetch all seekers
		var seekers []models.Seeker
		//log.Println("Fetching all seekers from the database...")
		if err := w.DB.Find(&seekers).Error; err != nil {
			log.Printf("Error fetching seekers: %v", err)
			time.Sleep(time.Minute) // Retry after a minute in case of an error
			continue
		}
		//log.Printf("Fetched %d seekers", len(seekers))

		// Loop through all seekers to calculate and store match scores
		for _, seeker := range seekers {
			//log.Printf("Processing Seeker ID: %s", seeker.AuthUserID)

			// Fetch user's preferred job titles
			var preferred models.PreferredJobTitle
			//log.Printf("Fetching preferred job titles for Seeker ID: %s", seeker.AuthUserID)
			if err := w.DB.Where("auth_user_id = ?", seeker.AuthUserID).First(&preferred).Error; err != nil {
				log.Printf("Error fetching preferred job titles for Seeker ID %s: %v", seeker.AuthUserID, err)
				continue
			}
			//log.Printf("Preferred job titles fetched: %+v", preferred)

			// Collect preferred titles
			var preferredTitles []string
			if preferred.PrimaryTitle != "" {
				preferredTitles = append(preferredTitles, preferred.PrimaryTitle)
			}
			if preferred.SecondaryTitle != nil && *preferred.SecondaryTitle != "" {
				preferredTitles = append(preferredTitles, *preferred.SecondaryTitle)
			}
			if preferred.TertiaryTitle != nil && *preferred.TertiaryTitle != "" {
				preferredTitles = append(preferredTitles, *preferred.TertiaryTitle)
			}

			if len(preferredTitles) == 0 {
				log.Printf("No preferred job titles set for Seeker ID %s", seeker.AuthUserID)
				continue
			}
			//log.Printf("Preferred job titles: %+v", preferredTitles)

			// Title filtering (case-insensitive)
			var conditions []string
			var values []interface{}
			for _, title := range preferredTitles {
				conditions = append(conditions, "LOWER(title) LIKE ?")
				values = append(values, "%"+strings.ToLower(title)+"%")
			}
			whereClause := strings.Join(conditions, " OR ")

			// Query LinkedIn jobs based on preferred titles
			//log.Println("Querying LinkedIn jobs...")
			var linkedinJobs []models.LinkedInJobMetaData
			if err := w.DB.Where(whereClause, values...).Find(&linkedinJobs).Error; err != nil {
				log.Printf("Error fetching LinkedIn jobs for Seeker ID %s: %v", seeker.AuthUserID, err)
				continue
			}
			//log.Printf("Fetched %d LinkedIn jobs for Seeker ID %s", len(linkedinJobs), seeker.AuthUserID)

			// Query Xing jobs based on preferred titles
			//log.Println("Querying Xing jobs...")
			var xingJobs []models.XingJobMetaData
			if err := w.DB.Where(whereClause, values...).Find(&xingJobs).Error; err != nil {
				//log.Printf("Error fetching Xing jobs for Seeker ID %s: %v", seeker.AuthUserID, err)
				continue
			}
			//log.Printf("Fetched %d Xing jobs for Seeker ID %s", len(xingJobs), seeker.AuthUserID)

			// Loop through LinkedIn jobs and calculate match scores
			for _, job := range linkedinJobs {
				// Check if the match score already exists for this job and seeker
				var existingMatchScore models.MatchScore
				if err := w.DB.Where("seeker_id = ? AND job_id = ?", seeker.AuthUserID, job.JobID).First(&existingMatchScore).Error; err == nil {
					// Match score already exists, skip this job
					log.Printf("Match score already exists for SeekerID: %s and JobID: %s, skipping", seeker.AuthUserID, job.JobID)
					continue
				}

				// If not, calculate and store the match score
				err := w.calculateAndStoreMatchScore(seeker.AuthUserID, job.JobID, "linkedin")
				if err != nil {
					log.Printf("Error calculating match score for Seeker ID %s and Job ID %s: %v", seeker.AuthUserID, job.JobID, err)
				} else {
					log.Printf("Match score calculated and stored for Seeker ID %s and Job ID %s", seeker.AuthUserID, job.JobID)
				}
			}

			// Loop through Xing jobs and calculate match scores
			for _, job := range xingJobs {
				// Check if the match score already exists for this job and seeker
				var existingMatchScore models.MatchScore
				if err := w.DB.Where("seeker_id = ? AND job_id = ?", seeker.AuthUserID, job.JobID).First(&existingMatchScore).Error; err == nil {
					// Match score already exists, skip this job
					log.Printf("Match score already exists for SeekerID: %s and JobID: %s, skipping", seeker.AuthUserID, job.JobID)
					continue
				}

				// If not, calculate and store the match score
				err := w.calculateAndStoreMatchScore(seeker.AuthUserID, job.JobID, "xing")
				if err != nil {
					log.Printf("Error calculating match score for Seeker ID %s and Job ID %s: %v", seeker.AuthUserID, job.JobID, err)
				} else {
					log.Printf("Match score calculated and stored for Seeker ID %s and Job ID %s", seeker.AuthUserID, job.JobID)
				}
			}
		}

		// Sleep for 1 hour after one full pass through seekers and jobs
		log.Println("MatchScoreWorker completed a full cycle, sleeping for 1 minute...")
		time.Sleep(time.Minute)
	}
}
