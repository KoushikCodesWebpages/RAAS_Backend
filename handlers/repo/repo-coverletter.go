package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"os"

	"RAAS/config" // Replace with the correct path to your config package
)

type CoverLetterInput struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	RecipientTitle  string `json:"recipient_title"`
	CompanyName     string `json:"company_name"`
	CompanyLocation string `json:"company_location"`
	Body            string `json:"body"`
	Closing         string `json:"closing"`
}


// GenerateCoverLetterDocx generates a cover letter document using input data
func GenerateCoverLetterDocx(input CoverLetterInput, config *config.Config) (string, error) {
	// Get the API URL and Key from the config
	apiURL := config.CL_Url
	apiKey := config.GEN_API_KEY

	// Check if the required fields are present
	if apiURL == "" || apiKey == "" {
		return "", fmt.Errorf("COVER_LETTER_API_URL or COVER_LETTER_API_KEY is missing in config")
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
	outFile := "cover_letter.docx"
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

