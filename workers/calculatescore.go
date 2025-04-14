package workers

import (
	
	"encoding/json"
	"fmt"
	// "io"
	// "bytes"
	// "net/http"
	"os"
	"github.com/joho/godotenv"
	"RAAS/models"
	//"log"
	"math"
	"errors"
	"strings"
)

// RoundRobinModelIndex keeps track of the current model index for round-robin
var (
	RoundRobinModelIndex int
)

// LoadHFModels loads the Hugging Face models from environment variables
func LoadHFModels() ([]string, error) {
	var models []string
	for i := 1; i <= 10; i++ {
		model := os.Getenv(fmt.Sprintf("HF_MODEL_FOR_MS_%d", i))
		if model == "" {
			return nil, fmt.Errorf("model HF_MODEL_%d is not defined", i)
		}
		models = append(models, model)
	}
	//log.Printf("Loaded Hugging Face models: %+v", models)  // Debugging log
	return models, nil
}

// CalculateMatchScore calculates the match score between the seeker and the job
func CalculateMatchScore(seeker models.Seeker, job interface{}) (float64, error) {
	if err := godotenv.Load(); err != nil {
		return 0, fmt.Errorf("error loading .env file: %v", err)
	}
	//log.Println("Loaded environment variables.")

	// Load Hugging Face API key
	hfAPIKey := os.Getenv("HF_API_KEY")
	if hfAPIKey == "" {
		return 0, fmt.Errorf("hugging Face API key not found")
	}

	// Load Hugging Face models
	modelsList, err := LoadHFModels()
	if err != nil {
		return 0, fmt.Errorf("error loading models: %v", err)
	}

	currentModel := modelsList[RoundRobinModelIndex]
	RoundRobinModelIndex = (RoundRobinModelIndex + 1) % len(modelsList)

	// === Professional Summary ===
	var summary models.ProfessionalSummary
	if err := models.DB.Where("auth_user_id = ?", seeker.AuthUserID).First(&summary).Error; err != nil {
		return 0, fmt.Errorf("failed to find professional summary for SeekerID %s: %v", seeker.AuthUserID, err)
	}

	var skills []string
	if err := json.Unmarshal(summary.Skills, &skills); err != nil {
		return 0, fmt.Errorf("failed to parse skills for SeekerID %s: %v", seeker.AuthUserID, err)
	}
	skillsStr := joinStrings(skills, ", ")

	// === Work Experience ===
	var workExperiences []models.WorkExperience
	if err := models.DB.Where("auth_user_id = ?", seeker.AuthUserID).Find(&workExperiences).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch work experience for SeekerID %s: %v", seeker.AuthUserID, err)
	}

	var workExpText string
	for _, we := range workExperiences {
		workExpText += fmt.Sprintf("Job Title: %s. Responsibilities: %s. ", we.JobTitle, we.KeyResponsibilities)
	}

	// === Education ===
	var educations []models.Education
	if err := models.DB.Where("auth_user_id = ?", seeker.AuthUserID).Find(&educations).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch education for SeekerID %s: %v", seeker.AuthUserID, err)
	}

	var eduText string
	for _, edu := range educations {
		eduText += fmt.Sprintf("Degree: %s in %s. Achievements: %s. ", edu.Degree, edu.FieldOfStudy, edu.Achievements)
	}

	// === Certificates ===
	var certificates []models.Certificate
	if err := models.DB.Where("auth_user_id = ?", seeker.AuthUserID).Find(&certificates).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch certificates for SeekerID %s: %v", seeker.AuthUserID, err)
	}

	var certText string
	for _, cert := range certificates {
		certText += fmt.Sprintf("Certificate: %s. ", cert.CertificateName)
	}

	// === Languages ===
	var languages []models.Language
	if err := models.DB.Where("auth_user_id = ?", seeker.AuthUserID).Find(&languages).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch languages for SeekerID %s: %v", seeker.AuthUserID, err)
	}

	var langText string
	for _, lang := range languages {
		langText += fmt.Sprintf("Language: %s (%s). ", lang.LanguageName, lang.ProficiencyLevel)
	}

	// === Preferred Job Titles ===
	var titles models.PreferredJobTitle
	if err := models.DB.Where("auth_user_id = ?", seeker.AuthUserID).First(&titles).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch preferred job titles for SeekerID %s: %v", seeker.AuthUserID, err)
	}

	jobTitles := "Preferred Job Titles: " + titles.PrimaryTitle
	if titles.SecondaryTitle != nil {
		jobTitles += ", " + *titles.SecondaryTitle
	}
	if titles.TertiaryTitle != nil {
		jobTitles += ", " + *titles.TertiaryTitle
	}
	jobTitles += "."

	// === Final Seeker Text ===
	seekerText := fmt.Sprintf(
		"Skills: %s. About: %s. Work Experience: %s Education: %s Certificates: %s Languages: %s %s",
		skillsStr, summary.About, workExpText, eduText, certText, langText, jobTitles,
	)

	// === Job Text from Metadata + Description ===
	var jobText string

	switch jobMeta := job.(type) {
	case models.LinkedInJobMetaData:
		var jobDesc models.LinkedInJobDescription
		if err := models.DB.Where("job_id = ?", jobMeta.JobID).First(&jobDesc).Error; err != nil {
			return 0, fmt.Errorf("failed to fetch LinkedIn job description for JobID %s: %v", jobMeta.JobID, err)
		}
		jobText = fmt.Sprintf("Title: %s. Description: %s. Skills: %s. Type: %s.",
			jobMeta.Title, jobDesc.JobDescription, jobDesc.Skills, jobDesc.JobType,
		)

	case models.XingJobMetaData:
		var jobDesc models.XingJobDescription
		if err := models.DB.Where("job_id = ?", jobMeta.JobID).First(&jobDesc).Error; err != nil {
			return 0, fmt.Errorf("failed to fetch Xing job description for JobID %s: %v", jobMeta.JobID, err)
		}
		jobText = fmt.Sprintf("Title: %s. Description: %s. Skills: %s. Type: %s.",
			jobMeta.Title, jobDesc.JobDescription, jobDesc.Skills, jobDesc.JobType,
		)

	default:
		return 0, fmt.Errorf("unsupported job type")
	}

	// === Hugging Face API Call ===
	//matchScore, err := getDirectMatchScoreFromHuggingFace(hfAPIKey, seekerText, jobText, currentModel)

	matchScore,err := CosineSimilarity(seekerText, jobText)
	if err != nil {
		return 0, fmt.Errorf("error getting match score from model %s: %v", currentModel, err)
	}

	return matchScore, nil
}
// Apply sigmoid function for scaling
func sigmoid(x float64) float64 {
	return 100 / (1 + math.Exp(-x)) // Converts into a range of 0 to 100
}

// CosineSimilarity calculates the cosine similarity between two text strings and applies sigmoid scaling
func CosineSimilarity(text1, text2 string) (float64, error) {
	// Tokenize the texts by splitting into words
	tokens1 := tokenize(text1)
	tokens2 := tokenize(text2)

	// If either of the texts results in no tokens, return an error
	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0, errors.New("one of the input texts is empty after tokenization")
	}

	// Calculate term frequencies (TF)
	tf1 := termFrequency(tokens1)
	tf2 := termFrequency(tokens2)

	// Calculate dot product and magnitudes
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	// Calculate the dot product and magnitudes
	for word, freq1 := range tf1 {
		freq2 := tf2[word]
		dotProduct += freq1 * freq2
		magnitude1 += freq1 * freq1
	}

	for _, freq2 := range tf2 {
		magnitude2 += freq2 * freq2
	}

	// Compute the cosine similarity
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0, errors.New("one of the vectors has zero magnitude")
	}

	// Calculate cosine similarity in the range of 0 to 1
	cosineSim := dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))

	// Apply sigmoid transformation to scale the score between 0 and 100
	matchScore := sigmoid(cosineSim * 10) // Amplifying cosine similarity a bit before sigmoid

	if matchScore < 0 {
		matchScore = 0
	} else if matchScore > 100 {
		matchScore = 100
	}

	return matchScore, nil
}



// tokenize splits a string into words (tokens)
func tokenize(text string) []string {
	// Convert the text to lowercase and split by non-alphanumeric characters
	// This basic tokenizer can be improved by using a more sophisticated library.
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

// termFrequency calculates the term frequency of each word in a list of tokens
func termFrequency(tokens []string) map[string]float64 {
	tf := make(map[string]float64)
	for _, token := range tokens {
		tf[token]++
	}
	// Normalize by the total number of words
	for word := range tf {
		tf[word] /= float64(len(tokens))
	}
	return tf
}



// func CosineSimilarity(text1, text2 string) (float64, error) {
// 	// Tokenize the texts by splitting into words
// 	tokens1 := tokenize(text1)
// 	tokens2 := tokenize(text2)

// 	// If either of the texts results in no tokens, return an error
// 	if len(tokens1) == 0 || len(tokens2) == 0 {
// 		return 0, errors.New("one of the input texts is empty after tokenization")
// 	}

// 	// Calculate term frequencies (TF)
// 	tf1 := termFrequency(tokens1)
// 	tf2 := termFrequency(tokens2)

// 	// Calculate dot product and magnitudes
// 	dotProduct := 0.0
// 	magnitude1 := 0.0
// 	magnitude2 := 0.0

// 	// Calculate the dot product and magnitudes
// 	for word, freq1 := range tf1 {
// 		freq2 := tf2[word]
// 		dotProduct += freq1 * freq2
// 		magnitude1 += freq1 * freq1
// 	}

// 	for _, freq2 := range tf2 {
// 		magnitude2 += freq2 * freq2
// 	}

// 	// Compute the cosine similarity
// 	if magnitude1 == 0 || magnitude2 == 0 {
// 		return 0, errors.New("one of the vectors has zero magnitude")
// 	}

// 	// Calculate cosine similarity in the range of 0 to 1
// 	cosineSim := dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))

// 	// Convert the cosine similarity to the range of 1 to 100
// 	matchScore := cosineSim * 100
// 	if matchScore < 0 {
// 		matchScore = 0 // Ensure score is non-negative
// 	} else if matchScore > 100 {
// 		matchScore = 100 // Cap the score at 100
// 	}

// 	return matchScore, nil
// }




// getDirectMatchScoreFromHuggingFace sends a request to Hugging Face API to get the match score
// func getDirectMatchScoreFromHuggingFace(apiKey, seekerText, jobText, model string) (float64, error) {
// 	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model)
// 	payload := map[string]interface{}{
// 		"inputs": fmt.Sprintf("%s %s", seekerText, jobText), // Combine both texts for matching
// 	}

// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return 0, fmt.Errorf("error marshaling payload: %v", err)
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return 0, fmt.Errorf("error creating request: %v", err)
// 	}

// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
// 	req.Header.Add("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return 0, fmt.Errorf("error making request to Hugging Face: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Check for successful response
// 	if resp.StatusCode != http.StatusOK {
// 		return 0, fmt.Errorf("error: received non-200 response code: %d", resp.StatusCode)
// 	}

// 	// Read and parse the response
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return 0, fmt.Errorf("error reading response: %v", err)
// 	}

// 	var response struct {
// 		Score float64 `json:"score"`
// 	}
// 	if err := json.Unmarshal(body, &response); err != nil {
// 		return 0, fmt.Errorf("error unmarshaling response: %v", err)
// 	}

// 	return response.Score, nil
// }


func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}

