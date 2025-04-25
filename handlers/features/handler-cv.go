package features

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"


// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"

// 	"RAAS/models"
// 	"RAAS/config"

// 	"time"
// 	"RAAS/handlers/repo"
// 	"log"
// 	// "bytes"
// 	// "log"
// )

// type CVData struct {
// 	JobID            string                 `json:"job_id"`
// 	Name             string                 `json:"name"`
// 	Designation      string                 `json:"designation"`
// 	Contact          string                 `json:"contact"`
// 	ProfileSummary   string                 `json:"profile_summary"`
// 	SkillsAndTools   []string               `json:"skills_and_tools"`
// 	Education        []EducationData        `json:"education"`
// 	ExperienceSummary []ExperienceSummaryData `json:"experience_summary"`
// 	Languages        []string               `json:"languages"`
// }


// type EducationData struct {
// 	Years       string   `json:"years"`
// 	Institution string   `json:"institution"`
// 	Details     []string `json:"details"`
// }

// type ExperienceSummaryData struct {
// 	Title  string   `json:"title"`
// 	Bullets []string `json:"bullets"`
// }



// // CVRequest struct to receive CV generation request
// type CVRequest struct {
// 	// Add any required fields here
// }

// // CVHandler struct
// type CVHandler struct {
// 	db     *gorm.DB
// }

// func NewCVHandler(db *gorm.DB, cfg *config.Config) *CVHandler {
// 	return &CVHandler{
// 		db:     db,
// 	}
// }
// func (h *CVHandler) PostCV(c *gin.Context) {

// 	var req CoverLetterAndCVRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid job_id in request body"})
// 		return
// 	}
// 	jobID := req.JobID
// 	// Step 1: Extract user information from JWT claims
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	// Step 2: Get Contact details


// 	var authUser models.AuthUser
// 	if err := h.db.Where("id = ?", userID).First(&authUser).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve auth user info"})
// 		return
// 	}

// 	//Email PHONE
// 	email := authUser.Email
// 	phone := authUser.Phone


// 	//Step 3: Get seeker model

// 	var seeker models.Seeker
// 	if err := h.db.Where("auth_user_id = ?", userID).First(&seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve seeker data"})
// 		return
// 	}

// 	if seeker.DailyGeneratableCV > 0 {
// 		seeker.DailyGeneratableCV -= 1
// 		if err := h.db.Save(&seeker).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seeker data"})
// 			return
// 		}
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Daily CV generation limit exceeded"})
// 		return
// 	}

//     // Step 3: Get personal info from seeker

// 	type PersonalInfo struct {
// 		FirstName      string `json:"firstName"`
// 		SecondName     string `json:"secondName"`
// 		Address        string `json:"address"`
// 		LinkedInProfile string `json:"linkedInProfile"`
// 		DateOfBirth    string `json:"dateOfBirth"`
// 	}
// 	var personalInfo PersonalInfo
// 	if err := json.Unmarshal(seeker.PersonalInfo, &personalInfo); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse personal information", "details": err.Error()})
// 		return
// 	}	

	
//     // Name and contact details
//     var fullName string
//     if personalInfo.SecondName != "" {
//         fullName = fmt.Sprintf("%s %s", personalInfo.FirstName, personalInfo.SecondName)
//     } else {
//         fullName = personalInfo.FirstName
//     }

//     // Address and LinkedIn profile
//     address := personalInfo.Address
//     linkedin := personalInfo.LinkedInProfile



// 	// Step 4: Get professional summary (includes skills)

// 	type ProfessionalSummary struct {
// 		About  string   `json:"about"`
// 		Skills []string `json:"skills"`
// 	}

// 	var summary ProfessionalSummary
// 	if err := json.Unmarshal(seeker.ProfessionalSummary, &summary); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse professional summary", "details": err.Error()})
// 		return
// 	}

// 	// Profile summary
// 	profile_summary := summary.About

// 	// Skills
// 	skills := summary.Skills

// // Step 6: Get education details
// type Education struct {
// 	Degree       string    `json:"degree"`
// 	StartDate    string    `json:"startDate"`
// 	EndDate      string    `json:"endDate"`
// 	Institution  string    `json:"institution"`
// 	Achievements string    `json:"achievements"`
// }

// var educations []Education
// if err := json.Unmarshal(seeker.Educations, &educations); err != nil {
// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse education details", "details": err.Error()})
// 	return
// }

// // Validate and parse start and end dates
// educationData := make([]struct {
// 	Years      string   `json:"years"`
// 	Institution string  `json:"institution"`
// 	Details    []string `json:"details"`
// }, 0)

// for _, edu := range educations {
// 	// Parse start date
// 	startDate, err := time.Parse("2006-01-02", edu.StartDate)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid start date format", "details": err.Error()})
// 		return
// 	}

// 	// Parse end date (if available)
// 	endDate := "Present"
// 	if edu.EndDate != "" {
// 		endDateParsed, err := time.Parse("2006-01-02", edu.EndDate)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid end date format", "details": err.Error()})
// 			return
// 		}
// 		endDate = endDateParsed.Format("2006")
// 	}

// 	// Prepare education details
// 	details := []string{edu.Degree}
// 	if edu.Achievements != "" {
// 		details = append(details, edu.Achievements)
// 	}

// 	educationData = append(educationData, struct {
// 		Years      string   `json:"years"`
// 		Institution string  `json:"institution"`
// 		Details    []string `json:"details"`
// 	}{
// 		Years:      fmt.Sprintf("%s - %s", startDate.Format("2006"), endDate),
// 		Institution: edu.Institution,
// 		Details:    details,
// 	})
// }

	
	
// 		// Step 7: Get work experience details
// 	type WorkExperience struct {
// 		JobTitle           string `json:"jobTitle"`
// 		CompanyName        string `json:"companyName"`
// 		EmploymentType     string `json:"employmentType"`
// 		StartDate          string `json:"startDate"`
// 		EndDate            string `json:"endDate"` // Required
// 		KeyResponsibilities string `json:"keyResponsibilities"`
// 	}

// 	var workExperiences []WorkExperience
// 	if err := json.Unmarshal(seeker.WorkExperiences, &workExperiences); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse work experience", "details": err.Error()})
// 		return
// 	}

// 	// Validate and parse work experience dates
// 	var workExpText string
// 	experienceSummaryData := make([]struct {
// 		Title   string   `json:"title"`
// 		Bullets []string `json:"bullets"`
// 	}, 0)

// 	for _, we := range workExperiences {
// 		// Parse start date
// 		startDate, err := time.Parse("2006-01-02", we.StartDate)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid start date format for work experience", "details": err.Error()})
// 			return
// 		}

// 		// Parse end date (if available)
// 		endDate := "Present"
// 		if we.EndDate != "" {
// 			endDateParsed, err := time.Parse("2006-01-02", we.EndDate)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid end date format for work experience", "details": err.Error()})
// 				return
// 			}
// 			endDate = endDateParsed.Format("2006")
// 		}

// 		// Append work experience summary text
// 		workExpText += fmt.Sprintf("Job Title: %s at %s. Responsibilities: %s. ", we.JobTitle, we.CompanyName, we.KeyResponsibilities)

// 		// Split key responsibilities into bullets
// 		bullets := strings.Split(we.KeyResponsibilities, ".")
// 		for i, bullet := range bullets {
// 			bullets[i] = strings.TrimSpace(bullet)
// 		}
// 		bullets = removeEmptyStrings(bullets) // Remove empty strings

// 		// Add to experience summary data
// 		experienceSummaryData = append(experienceSummaryData, struct {
// 			Title   string   `json:"title"`
// 			Bullets []string `json:"bullets"`
// 		}{
// 			Title:   fmt.Sprintf("%s at %s (%s - %s)", we.JobTitle, we.CompanyName, startDate.Format("2006"), endDate),
// 			Bullets: bullets,
// 		})
// 	}


// 	// Step 8: Get languages details
// 	type Language struct {
// 		LanguageName     string `json:"language"`
// 		ProficiencyLevel string `json:"proficiency"`
// 	}

// 	var languages []Language
// 	if err := json.Unmarshal(seeker.Languages, &languages); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse language details", "details": err.Error()})
// 		return
// 	}

// 	// Format the language data for the CV
// 	languageData := make([]string, 0)
// 	for _, lang := range languages {
// 		// Print every language's name and proficiency level for debugging
// 		fmt.Printf("Language Name: %s, Proficiency Level: %s\n", lang.LanguageName, lang.ProficiencyLevel)
		
// 		languageData = append(languageData, fmt.Sprintf("%s - %s", lang.LanguageName, lang.ProficiencyLevel))
// 	}

// 	// Print the formatted language data to the console for debugging
// 	fmt.Println("Formatted Language Data:", languageData)


// 	// Step 9 Get Job Title

// 	var job models.Job
// 	if err := h.db.Where("job_id = ?", jobID).First(&job).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job metadata"})
// 		return
// 	}
	
// 	// Job details are now in jobMetaData
// 	jobTitle := job.Title

// 	contact := fmt.Sprintf("Email: %s\nPhone: %s\nAddress: %s\nLinkedIn: %s", email, phone, address, linkedin)

// 	cvInput := repo.CVInput{
//         Name:             fullName,
//         Designation:      jobTitle,
//         Contact:          contact,
//         ProfileSummary:   profile_summary,
//         SkillsAndTools:   skills,
//         Education:        educationData,
//         ExperienceSummary: experienceSummaryData,
//         Languages:        languageData,
//     }
    
//     // Generate CV (docx)
//     cvData, err := repo.GenerateCVDocx(cvInput)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CV document", "details": err.Error()})
//         return
//     }

//     // Initialize the media upload handler with the Azure Blob Storage service client
//     mediaUploadHandler := NewMediaUploadHandler(GetBlobServiceClient())

//     // Upload the CV document (cvData) directly to Azure Blob Storage
//     cvFileURL, err := mediaUploadHandler.UploadGeneratedFile(c, "cv-container", "cv.docx", cvData)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload CV to Azure", "details": err.Error()})
//         return
//     }
//     // Update the job application record to mark CV as generated
//     result := h.db.Model(&models.SelectedJobApplication{}).Where("auth_user_id = ? AND job_id = ?", userID, jobID).Update("cv_generated", true)
//     if result.Error != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update CvGenerated field"})
//         return
//     }

//     // Now, send the CV file as an attachment
//     c.Header("Content-Disposition", "attachment; filename=cv.docx")
//     c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", cvData)

//     // Optionally, you can also update the database model with cv_url (in case you want to store the link)
//     cvRecord := models.CV{
//         AuthUserID: userID,
//         JobID:      jobID,
//         CVUrl:      cvFileURL, // URL of the uploaded CV
//     }

//     if err := h.db.Create(&cvRecord).Error; err != nil {
//         log.Printf("[ERROR] Failed to create CV record: %v", err)
//     }
// }

// func (h *CVHandler) GetCV(c *gin.Context) {
//     // Retrieve the user ID from the context
//     userID := c.MustGet("userID").(uuid.UUID)

//     // Retrieve the Job ID from the query parameters (or you could use c.PostForm if it's a POST request)
//     jobID := c.DefaultQuery("jobID", "") // If it's a query parameter, use c.DefaultQuery to get the jobID

//     if jobID == "" {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
//         return
//     }

//     var cv models.CV

//     // Query the CV table for the specified AuthUserID and JobID
//     if err := h.db.Where("auth_user_id = ? AND job_id = ?", userID, jobID).First(&cv).Error; err != nil {
//         // Return error if no CV is found
//         c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
//         return
//     }

//     // Check if the CV URL exists
//     if cv.CVUrl == "" {
//         c.JSON(http.StatusNotFound, gin.H{"error": "CV file URL not found"})
//         return
//     }

//     // Attempt to download the CV from the provided URL
//     fileURL := cv.CVUrl
//     response, err := http.Get(fileURL)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download the file", "details": err.Error()})
//         return
//     }
//     defer response.Body.Close()

//     // Check if the file was successfully fetched
//     if response.StatusCode != http.StatusOK {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download the file, received status: " + response.Status})
//         return
//     }

//     // Set headers to indicate a file download response
//     c.Header("Content-Disposition", "attachment; filename=cv.docx")
    
//     // Custom headers map (optional)
//     headers := map[string]string{
//         "Content-Type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
//     }

//     // Send the file to the user with the required headers
//     c.DataFromReader(http.StatusOK, response.ContentLength, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", response.Body, headers)
// }


// func removeEmptyStrings(s []string) []string {
// 	var r []string
// 	for _, str := range s {
// 		if str != "" {
// 			r = append(r, str)
// 		}
// 	}
// 	return r
// }


