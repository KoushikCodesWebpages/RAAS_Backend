package handlers

import (
	"RAAS/dto"
	"RAAS/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// MatchScorePOST computes and stores match score for a job
func MatchScorePOST(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var req dto.MatchScoreRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_id is required", "details": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	// Check if score already exists
	var existing models.JobMatchScore
	if err := db.Where("auth_user_id = ? AND job_id = ?", userID, req.JobID).First(&existing).Error; err == nil {
		response := dto.MatchScoreResponse{
			Score:     existing.Score,
			Source:    "cached",
			JobID:     existing.JobID,
			UserID:    existing.AuthUserID,
			Platform:  existing.Platform,
			MatchedAt: existing.MatchedAt,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// Calculate new match score
	score := calculateMatchScore(req.JobID, userID)
	match := models.JobMatchScore{
		AuthUserID: userID,
		JobID:      req.JobID,
		Platform:   detectPlatform(req.JobID),
		Score:      score,
		MatchedAt:  time.Now(),
	}

	if err := db.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store match score", "details": err.Error()})
		return
	}

	response := dto.MatchScoreResponse{
		Score:     match.Score,
		Source:    "calculated",
		JobID:     match.JobID,
		UserID:    match.AuthUserID,
		Platform:  match.Platform,
		MatchedAt: match.MatchedAt,
	}

	c.JSON(http.StatusOK, response)
}

// MatchScoreGET retrieves cached match score for a job
func MatchScoreGET(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	jobID := c.Query("job_id")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_id is required"})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	var score models.JobMatchScore
	if err := db.Where("auth_user_id = ? AND job_id = ?", userID, jobID).First(&score).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"matched": false,
			"message": "This job has not been processed for a match yet.",
			"job_id":  jobID,
		})
		return
	}

	response := dto.MatchScoreResponse{
		Score:     score.Score,
		Source:    "cached",
		JobID:     score.JobID,
		UserID:    score.AuthUserID,
		Platform:  score.Platform,
		MatchedAt: score.MatchedAt,
	}
	c.JSON(http.StatusOK, response)
}

// detectPlatform determines the job platform from job ID prefix
func detectPlatform(jobID string) string {
	switch {
	case strings.HasPrefix(jobID, "L"):
		return "linkedin"
	case strings.HasPrefix(jobID, "X"):
		return "xing"
	default:
		return "unknown"
	}
}

// calculateMatchScore simulates match score calculation
func calculateMatchScore(jobID string, userID uuid.UUID) float64 {
	return math.Round(rand.Float64()*10000) / 100 // e.g., 87.43
}
