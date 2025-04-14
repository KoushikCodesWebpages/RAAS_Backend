package repo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	

	"io"
	"os"
	"RAAS/config"
	// Assuming your models are stored here
)

type RequestData struct {
	Inputs string `json:"inputs"`
}

type ResponseData struct {
	GeneratedText string `json:"generated_text"`
}

type CVInput struct {
	Name             string   `json:"name"`
	Designation      string   `json:"designation"`
	Contact          string   `json:"contact"`
	ProfileSummary   string   `json:"profile_summary"`
	SkillsAndTools   []string `json:"skills_and_tools"`
	Education        []struct {
		Years      string   `json:"years"`
		Institution string `json:"institution"`
		Details    []string `json:"details"`
	} `json:"education"`
	ExperienceSummary []struct {
		Title  string   `json:"title"`
		Bullets []string `json:"bullets"`
	} `json:"experience_summary"`
	Languages []string `json:"languages"`
}

// GenerateCVDocx generates a CV document using input data
func GenerateCVDocx(input CVInput, config *config.Config) (string, error) {
	// Get the API URL and Key from the config
	apiURL := config.CV_Url
	apiKey := config.GEN_API_KEY

	// Check if the required fields are present
	if apiURL == "" || apiKey == "" {
		return "", fmt.Errorf("CV_API_URL or COVER_CV_API_KEY is missing in config")
	}

	// Marshal the input data into JSON
	jsonData, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("error marshaling input data: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Save the response body to a .docx file
	outFile := "resume.docx"
	out, err := os.Create(outFile)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error copying response to file: %v", err)
	}

	return outFile, nil
}


