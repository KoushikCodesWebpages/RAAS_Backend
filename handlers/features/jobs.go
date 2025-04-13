package features

import (
	"RAAS/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobDataHandler struct {
	db *gorm.DB
}

func NewJobDataHandler(db *gorm.DB) *JobDataHandler {
	return &JobDataHandler{db: db}
}

// LinkedIn Jobs
func (h *JobDataHandler) GetLinkedInJobs(c *gin.Context) {
	var jobs []models.LinkedInJobMetaData
	h.db.Find(&jobs)
	c.JSON(http.StatusOK, gin.H{"linkedin_jobs": jobs})
}

// Xing Jobs
func (h *JobDataHandler) GetXingJobs(c *gin.Context) {
	var jobs []models.XingJobMetaData
	h.db.Find(&jobs)
	c.JSON(http.StatusOK, gin.H{"xing_jobs": jobs})
}

// LinkedIn Application Links
func (h *JobDataHandler) GetLinkedInLinks(c *gin.Context) {
	var links []models.LinkedInJobApplicationLink
	h.db.Find(&links)
	c.JSON(http.StatusOK, gin.H{"linkedin_links": links})
}

// Xing Application Links
func (h *JobDataHandler) GetXingLinks(c *gin.Context) {
	var links []models.XingJobApplicationLink
	h.db.Find(&links)
	c.JSON(http.StatusOK, gin.H{"xing_links": links})
}

// LinkedIn Descriptions
func (h *JobDataHandler) GetLinkedInDescriptions(c *gin.Context) {
	var descs []models.LinkedInJobDescription
	h.db.Find(&descs)
	c.JSON(http.StatusOK, gin.H{"linkedin_descriptions": descs})
}

// Xing Descriptions
func (h *JobDataHandler) GetXingDescriptions(c *gin.Context) {
	var descs []models.XingJobDescription
	h.db.Find(&descs)
	c.JSON(http.StatusOK, gin.H{"xing_descriptions": descs})
}
