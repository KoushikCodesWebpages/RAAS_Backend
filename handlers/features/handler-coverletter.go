package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"bytes"
	"strings"

	"RAAS/models"
	"RAAS/config"
	"RAAS/handlers/repo"
)

// CoverLetterRequest struct to receive JobID
type CoverLetterRequest struct {
	JobID string `json:"job_id" binding:"required"`
}

// CoverLetterHandler struct
type CoverLetterHandler struct {
	db     *gorm.DB
}

func NewCoverLetterHandler(db *gorm.DB, cfg *config.Config) *CoverLetterHandler {
	return &CoverLetterHandler{
		db:     db,
	}
}



// GenerateCoverLetterBody function to generate a cover letter body based on user and job data
func (h *CoverLetterHandler) PostCoverLetter(c *gin.Context) {
	// Step 1: Parse job_id from request body
	var req CoverLetterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid job_id in request body"})
		return
	}
	jobID := req.JobID

	// Step 2: Extract user information from JWT claims
	userID := c.MustGet("userID").(uuid.UUID)

	// Step 3: Get user personal info
	var personalInfo models.PersonalInfo
	if err := h.db.Where("auth_user_id = ?", userID).First(&personalInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user personal information"})
		return
	}


	// Step 3b: Get auth user info (email, phone)
	var authUser models.AuthUser
	if err := h.db.Where("id = ?", userID).First(&authUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve auth user info"})
		return
	}


	// Step 4: Get professional summary (includes skills)
	var summary models.ProfessionalSummary
	if err := h.db.Where("auth_user_id = ?", userID).First(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve professional summary"})
		return
	}

	var skills []string
	if err := json.Unmarshal(summary.Skills, &skills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse skills"})
		return
	}
	skillsStr := joinStrings(skills, ", ")

	// Step 5: Get work experience
	var workExperiences []models.WorkExperience
	if err := h.db.Where("auth_user_id = ?", userID).Find(&workExperiences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve work experience"})
		return
	}
	var workExpText string
	for _, we := range workExperiences {
		workExpText += fmt.Sprintf("Job Title: %s at %s. Responsibilities: %s. ", we.JobTitle, we.CompanyName, we.KeyResponsibilities)
	}

	// Step 6: Get job metadata (LinkedIn or Xing)
	var companyName, jobTitle string
	var jobMetaData models.LinkedInJobMetaData
	if err := h.db.Where("job_id = ?", jobID).First(&jobMetaData).Error; err == nil {
		companyName = jobMetaData.Company
		jobTitle = jobMetaData.Title
	} else {
		var xingJobMetaData models.XingJobMetaData
		if err := h.db.Where("job_id = ?", jobID).First(&xingJobMetaData).Error; err == nil {
			companyName = xingJobMetaData.Company
			jobTitle = xingJobMetaData.Title
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job metadata"})
			return
		}
	}

	// Step 7: Get education
	var educations []models.Education
	if err := h.db.Where("auth_user_id = ?", userID).Find(&educations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve education details"})
		return
	}
	var eduText string
	for _, edu := range educations {
		eduText += fmt.Sprintf("Degree: %s from %s in %s. Achievements: %s. ", edu.Degree, edu.Institution, edu.FieldOfStudy, edu.Achievements)
	}

	// Step 8: Construct full name
	var fullName string
	if personalInfo.SecondName != nil {
		fullName = fmt.Sprintf("%s %s", personalInfo.FirstName, *personalInfo.SecondName)
	} else {
		fullName = personalInfo.FirstName
	}





	// Step 10: Generate cover letter body using Hugging Face models and input data
	coverLetterBody := generateCoverLetterBody(
		fullName,
		eduText,
		workExpText,
		skillsStr,
		companyName,
		jobTitle,
	)

	email := authUser.Email
	phone := authUser.Phone
	address := personalInfo.Address


	docInput := repo.CoverLetterInput{
		Name:            fullName,
		Email:           email,
		Phone:           phone,
		Address:         address,
		RecipientTitle:  fmt.Sprintf("Hiring Manager at %s", companyName),
		CompanyName:     companyName,
		CompanyLocation: "Germany", // Or fetch this from jobMetaData if available
		Body:            coverLetterBody,
		Closing:          "Sincerely",
	}

	docData, err := repo.GenerateCoverLetterDocx(docInput, config.Cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate cover letter document", "details": err.Error()})
		return
	}
	
	c.Header("Content-Disposition", "attachment; filename=cover_letter.docx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", []byte(docData))

}

// Helper function to join strings with commas
func joinStrings(arr []string, delimiter string) string {
	return fmt.Sprintf("%s", arr)
}

// Global variables
var (
	RoundRobinModelIndex int
	HFModels             []string
	APIKey               string
	HFBaseAPIURL         string
)

func LoadHFModels() ([]string, error) {
	models := []string{
		config.Cfg.HFModelForCL1,
		config.Cfg.HFModelForCL2,
		config.Cfg.HFModelForCL3,
		config.Cfg.HFModelForCL4,
		config.Cfg.HFModelForCL5,
		config.Cfg.HFModelForCL6,
		config.Cfg.HFModelForCL7,
		config.Cfg.HFModelForCL8,
		config.Cfg.HFModelForCL9,
		config.Cfg.HFModelForCL10,
	}

	for _, model := range models {
		if model == "" {
			log.Printf("Error loading Hugging Face models: one or more models are not defined")
			return nil, fmt.Errorf("one or more models are not defined")
		}
	}

	log.Printf("Loaded Hugging Face models: %+v", models)

	return models, nil
}


// prepareCoverLetterPrompt prepares the input prompt and API URL for the Hugging Face API
func prepareCoverLetterPrompt(name, education, experience, skills, company, role string) (string, string, error) {
	if name == "" || education == "" || experience == "" || skills == "" || company == "" || role == "" {
		log.Println("Error: Missing required field. All input fields must be provided.")
		return "", "", fmt.Errorf("error: Missing required information")
	}

	HFModels, err := LoadHFModels()
	if err != nil {
		log.Printf("Failed to load models: %v", err)
		return "", "", fmt.Errorf("error loading Hugging Face models: %w", err)
	}

	if len(HFModels) == 0 {
		log.Println("Error: HFModels is empty. Cannot generate cover letter.")
		return "", "", fmt.Errorf("error: No Hugging Face models available")
	}

	log.Printf("Loaded Hugging Face models: %v", HFModels)

	modelToUse := HFModels[RoundRobinModelIndex]
	log.Printf("Selected model: %s", modelToUse)

	if config.Cfg.HFBaseAPIUrl == "" {
		log.Println("Error: Hugging Face base API URL is not set")
		return "", "", fmt.Errorf("error: Hugging Face base API URL is not set")
	}
	apiURL := fmt.Sprintf("%s/%s", config.Cfg.HFBaseAPIUrl, modelToUse)

	RoundRobinModelIndex = (RoundRobinModelIndex + 1) % len(HFModels)
	prompt := fmt.Sprintf(`
	Write the body of a professional cover letter for %s applying for the %s role at %s.
	Education: %s. Experience: %s. Skills: %s.
	Write at least 3 well-structured paragraphs focusing on motivation, qualifications, and alignment with the role.
	Do not include greeting or closing.
	`, name, role, company, education, experience, skills)

	log.Printf("Generated prompt: %s", prompt)

	return prompt, apiURL, nil
}

// callHuggingFaceAPI sends a request to the Hugging Face API and returns the generated text
func callHuggingFaceAPI(prompt, apiURL string) (string, error) {
    requestBody, err := json.Marshal(map[string]interface{}{
        "inputs": prompt,
        "parameters": map[string]interface{}{
            "max_length": 1000,    // Increase the length of the response
        },
    })
    if err != nil {
        log.Printf("Error creating JSON request: %v", err)
        return "", fmt.Errorf("error: Failed to create request")
    }

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
    if err != nil {
        log.Printf("Error creating request: %v", err)
        return "", fmt.Errorf("error: Failed to create API request")
    }

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Cfg.HFAPIKey))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making API request: %v", err)
        return "", fmt.Errorf("error: Failed to reach AI service")
    }
    defer resp.Body.Close()

    log.Printf("API Response Status: %s", resp.Status)

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response: %v", err)
        return "", fmt.Errorf("error: Failed to read AI response")
    }

    log.Printf("Raw API Response: %s", body)

    var response []struct {
        GeneratedText string `json:"generated_text"`
    }

    if err := json.Unmarshal(body, &response); err != nil {
        log.Printf("Error parsing JSON response: %v", err)
        return "", fmt.Errorf("error: Failed to parse AI response")
    }

    if len(response) > 0 && response[0].GeneratedText != "" {
        log.Printf("Generated Text: %s", response[0].GeneratedText)
        return strings.TrimPrefix(response[0].GeneratedText, prompt), nil
    }

    log.Println("Error: No text generated.")
    return "", fmt.Errorf("error: No text generated")
}


func generateCoverLetterBody(name, education, experience, skills, company, role string) string {
	prompt, apiURL, err := prepareCoverLetterPrompt(name, education, experience, skills, company, role)
	if err != nil {
		log.Println("Error:", err)
		return err.Error()
	}

	coverLetter, err := callHuggingFaceAPI(prompt, apiURL)
	if err != nil {
		log.Println("Error:", err)
		return err.Error()
	}

	return coverLetter
}
