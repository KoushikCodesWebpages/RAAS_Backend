package features

import (
	"RAAS/handlers"
	"RAAS/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	// "time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoverLetterAndCVRequest struct {
	JobID string `json:"job_id" binding:"required"`
}

type CoverLetterHandler struct{}

func NewCoverLetterHandler() *CoverLetterHandler {
	return &CoverLetterHandler{}
}

func (h *CoverLetterHandler) PostCoverLetter(c *gin.Context) {
	// Get the database and collections
	db := c.MustGet("db").(*mongo.Database)
	jobCollection := db.Collection("jobs")
	seekerCollection := db.Collection("seekers")
	AuthUserCollection := db.Collection("auth_users")

	// Get the authenticated user's ID
	userID := c.MustGet("userID").(string)

	// Input: Expect a JobID to retrieve job details
	var input struct {
		JobID string `json:"job_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	// Fetch the job
	var job models.Job
	if err := jobCollection.FindOne(c, bson.M{"job_id": input.JobID}).Decode(&job); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
		return
	}

	// Fetch the seeker
	var seeker models.Seeker
	if err := seekerCollection.FindOne(c, bson.M{"auth_user_id": userID}).Decode(&seeker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching seeker data",
		})
		return
	}

	var authuser models.AuthUser
	if err := AuthUserCollection.FindOne(c, bson.M{"auth_user_id": userID}).Decode(&authuser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching authuser data",
		})
		return
	}

	// Get various details from seeker
	personalInfo, _ := handlers.GetPersonalInfo(&seeker)
	professionalSummary, _ := handlers.GetProfessionalSummary(&seeker)
	workExperience, _ := handlers.GetWorkExperience(&seeker)
	education, _ := handlers.GetEducation(&seeker)
	certificates, _ := handlers.GetCertificates(&seeker)
	languages, _ := handlers.GetLanguages(&seeker)

	// Construct the request body to match the API's expected format
	apiRequestData := map[string]interface{}{
		"user_details": map[string]interface{}{
			"firstname":        personalInfo.FirstName,
			"designation": seeker.PrimaryTitle,
			"address":     personalInfo.Address,
			"contact":     authuser.Phone,
			"email":       authuser.Email,
			// "portfolio":   personalInfo.Portfolio,
			"linkedin":    personalInfo.LinkedInProfile,
			// "tools":       seeker.Tools, // Assuming Tools field is populated correctly
			"skills":      professionalSummary.Skills, // Assuming Skills field is populated correctly
			"education":   education,
			"experience_summary": workExperience,
			"certifications": certificates,
			"languages":    languages,
		},
		"job_description": map[string]interface{}{
			"job_title":    job.Title,
			"company":      job.Company,
			"location":     job.Location,
			"job_type":     job.JobType, // Assuming JobType exists in the Job model
			"description": job.JobDescription, // Assuming Responsibilities exists in the Job model
			// "qualifications": job.Qualifications, // Assuming Qualifications exists in the Job model
			"skills":        job.Skills,
			// "benefits":      job.Benefits, // Assuming Benefits exists in the Job model
		},
	}

	// Call the external API to generate the cover letter DOCX
	docxContent, err := h.generateCoverLetter(apiRequestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error generating cover letter",
		})
		return
	}

	// Optionally, save the DOCX file or return it to the user
	// For now, sending the DOCX content back as response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", docxContent)
}

// Helper function to send POST request to the external cover letter generation API
func (h *CoverLetterHandler) generateCoverLetter(apiRequestData map[string]interface{}) ([]byte, error) {
	// Load environment variables
	apiURL := os.Getenv("COVER_LETTER_API_URL")
	apiKey := os.Getenv("COVER_CV_API_KEY")

	// Marshal cover letter data to JSON
	jsonData, err := json.Marshal(apiRequestData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling cover letter data: %v", err)
	}

	// Create a POST request to the API
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request to the cover letter API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error response from API: %v", string(body))
	}

	// Read the DOCX content from the response
	docxFileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading DOCX content: %v", err)
	}

	// Return the DOCX file content
	return docxFileContent, nil
}
