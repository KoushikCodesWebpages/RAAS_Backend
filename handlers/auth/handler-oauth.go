package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/gorm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"RAAS/config"
	"RAAS/models"
	"RAAS/security"
)

// Google OAuth Config
func getGoogleOauth2Config() oauth2.Config {
	// Print the config.Cfg to debug its values before using it
	log.Printf("Config.Cfg: %+v", config.Cfg)

	// Google OAuth2 Configuration
	return oauth2.Config{
		ClientID:     config.Cfg.Cloud.GoogleClientId,
		ClientSecret: config.Cfg.Cloud.GoogleClientSecret,
		RedirectURL:  config.Cfg.Cloud.GoogleRedirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

// GoogleLoginHandler handles the Google login OAuth flow
func GoogleLoginHandler(c *gin.Context) {
	// Step 1: Generate OAuth URL to redirect user to Google
	googleOauth2Config := getGoogleOauth2Config()
	authURL := googleOauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	log.Printf("Redirecting user to: %s", authURL)
	c.Redirect(http.StatusFound, authURL)
}

// GoogleCallbackHandler handles the Google OAuth callback
func GoogleCallbackHandler(c *gin.Context) {
	log.Println("GoogleCallbackHandler started")

	// Step 1: Get the authorization code from the query parameters
	code := c.DefaultQuery("code", "")
	if code == "" {
		log.Println("Authorization code is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is missing"})
		return
	}

	// Step 2: Exchange the code for an access token
	googleOauth2Config := getGoogleOauth2Config()
	log.Println("Exchanging authorization code for access token")
	token, err := googleOauth2Config.Exchange(c, code)
	if err != nil {
		log.Printf("Error exchanging token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token", "details": err.Error()})
		return
	}

	// Step 3: Get the user's profile info from Google
	log.Println("Fetching user info from Google")
	client := googleOauth2Config.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Printf("Error fetching user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info from Google", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Decode the user info from Google
	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Sub   string `json:"sub"` // This is the unique user ID from Google
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("Error decoding user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info", "details": err.Error()})
		return
	}

	log.Printf("User Info: Email=%s, Name=%s, GoogleID=%s", userInfo.Email, userInfo.Name, userInfo.Sub)

	// Step 4: Check if the user already exists in the database
	log.Println("Checking if user exists in the database")
	db := c.MustGet("db").(*gorm.DB)

	var user models.AuthUser
	if err := db.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Database error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
			return
		}

		// If user doesn't exist, create a new user with Google provider
		log.Println("Creating new user with Google provider")
		user = models.AuthUser{
			ID:            uuid.New(),
			Email:         userInfo.Email,
			Phone:         "", // Optional: you can ask for a phone number later
			Role:          "seeker", // Default role for new users
			Provider:      "google",
			EmailVerified: true, // Assume Google email is verified
			IsActive:      true,
			CreatedBy:     uuid.Nil,
			UpdatedBy:     uuid.Nil,
		}

		// Save AuthUser to DB
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Error creating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
			return
		}

		// Create Seeker profile for new user
		log.Println("Creating Seeker profile for new user")
		seeker := models.Seeker{
			AuthUserID:               user.ID,
			SubscriptionTier:         "free", // Default value for subscription tier
			DailySelectableJobsCount: 5,     // Default value
			DailyGeneratableCV:       100,    // Default value
			DailyGeneratableCoverletter: 100, // Default value
			TotalApplications:        0,      // Default value
			PersonalInfo:             nil,    // or initialize with an empty JSON object
			ProfessionalSummary:      nil,    // or initialize with an empty JSON object
			WorkExperiences:          nil,    // or initialize with an empty JSON object
			PrimaryTitle:             "",     // Empty initially
			SecondaryTitle:           nil,    // Empty initially
			TertiaryTitle:            nil,    // Empty initially
		}

		// Save Seeker to DB
		if err := db.Create(&seeker).Error; err != nil {
			log.Printf("Error creating seeker profile: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seeker profile", "details": err.Error()})
			return
		}
		log.Println("New user and seeker profile created successfully")
	}

	// Step 5: Issue JWT token for authenticated user
	log.Println("Issuing JWT token for authenticated user")
	tokenString, err := security.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT", "details": err.Error()})
		return
	}

	// Step 6: Respond with the JWT token
	log.Printf("JWT token generated successfully for user %s", user.Email)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
