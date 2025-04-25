package dataentry

// import (

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// 	"net/http"

// 	"RAAS/models"
// 	"RAAS/dto"
// 	"RAAS/handlers"
// )

// // PersonalInfoHandler struct
// type PersonalInfoHandler struct {
// 	DB *gorm.DB
// }

// // NewPersonalInfoHandler creates a new handler instance
// func NewPersonalInfoHandler(db *gorm.DB) *PersonalInfoHandler {
// 	return &PersonalInfoHandler{DB: db}
// }



// // CreatePersonalInfo handles creation of personal info within Seeker model
// func (h *PersonalInfoHandler) CreatePersonalInfo(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	// Bind input JSON to request struct
// 	var input dto.PersonalInfoRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
// 		return
// 	}

// 	// Find Seeker and handle error
// 	seeker, err := handlers.FindSeekerByUserID(h.DB, userID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	// Check if PersonalInfo is already filled
// 	if isFilled, err := handlers.IsFieldFilled(seeker.PersonalInfo); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check personal info", "details": err.Error()})
// 		return
// 	} else if isFilled {
// 		c.JSON(http.StatusConflict, gin.H{"error": "Personal info already filled"})
// 		return
// 	}

// 	// Validate DateOfBirth
// 	if input.DateOfBirth == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Date of birth cannot be empty"})
// 		return
// 	}

// 	// Set PersonalInfo and handle error
// 	if err := handlers.SetPersonalInfo(seeker, &input); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set personal info", "details": err.Error()})
// 		return
// 	}

// 	// Save updated Seeker and handle error
// 	if err := h.DB.Save(seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
// 		return
// 	}

// 	// Update the UserEntryTimeline - set PreferredJobTitlesCompleted to true
// 	var timeline models.UserEntryTimeline
// 	if err := h.DB.First(&timeline, "auth_user_id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User entry timeline not found"})
// 		return
// 	}
// 	// Return the PersonalInfoResponse
// 	c.JSON(http.StatusCreated, dto.PersonalInfoResponse{
// 		AuthUserID:      seeker.AuthUserID,
// 		FirstName:       input.FirstName,
// 		SecondName:      input.SecondName,
// 		DateOfBirth:     input.DateOfBirth,
// 		Address:         input.Address,
// 		LinkedInProfile: input.LinkedInProfile,
// 	})
// }







// // GetPersonalInfo retrieves personal info of the authenticated user
// func (h *PersonalInfoHandler) GetPersonalInfo(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	// Find Seeker and handle error
// 	seeker, err := handlers.FindSeekerByUserID(h.DB, userID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	// Check if PersonalInfo is empty or "null"
// 	if len(seeker.PersonalInfo) == 0 || string(seeker.PersonalInfo) == "null" {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Personal info not yet filled"})
// 		return
// 	}

// 	// Retrieve and return PersonalInfo
// 	personalInfo, err := handlers.GetPersonalInfo(seeker)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse personal info", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, personalInfo)
// }







// // UpdatePersonalInfo handles the update of personal info for Seeker model
// func (h *PersonalInfoHandler) UpdatePersonalInfo(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	var input dto.PersonalInfoRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
// 		return
// 	}

// 	// Find Seeker for the given user
// 	seeker, err := handlers.FindSeekerByUserID(h.DB, userID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	// Use SetPersonalInfo utility function to update
// 	if err := handlers.SetPersonalInfo(seeker, &input); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
// 		return
// 	}

// 	// Save the updated Seeker record
// 	if err := h.DB.Save(seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated personal info", "details": err.Error()})
// 		return
// 	}

// 	// Return the updated PersonalInfoResponse
// 	c.JSON(http.StatusOK, dto.PersonalInfoResponse{
// 		AuthUserID:      seeker.AuthUserID,
// 		FirstName:       input.FirstName,
// 		SecondName:      input.SecondName,
// 		DateOfBirth:     input.DateOfBirth,
// 		Address:         input.Address,
// 		LinkedInProfile: input.LinkedInProfile,
// 	})
// }






// // PatchPersonalInfo handles patching personal info for Seeker model
// func (h *PersonalInfoHandler) PatchPersonalInfo(c *gin.Context) {
// 	userID := c.MustGet("userID").(uuid.UUID)

// 	var input dto.PersonalInfoRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
// 		return
// 	}

// 	// Find Seeker for the given user
// 	seeker, err := handlers.FindSeekerByUserID(h.DB, userID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Seeker not found"})
// 		return
// 	}

// 	// Get current personal info and check for updates
// 	personalInfo, err := handlers.GetPersonalInfo(seeker)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse personal info", "details": err.Error()})
// 		return
// 	}

// 	// Reject attempts to update restricted fields
// 	if input.DateOfBirth != "" || input.Address != "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Updating DateOfBirth or Address is not allowed"})
// 		return
// 	}

// 	// Update only allowed fields
// 	if input.FirstName != "" {
// 		personalInfo.FirstName = input.FirstName
// 	}
// 	if input.SecondName != nil {
// 		personalInfo.SecondName = input.SecondName
// 	}
// 	if input.LinkedInProfile != nil {
// 		personalInfo.LinkedInProfile = input.LinkedInProfile
// 	}

// 	// Set updated personal info
// 	if err := handlers.SetPersonalInfo(seeker, personalInfo); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info", "details": err.Error()})
// 		return
// 	}

// 	// Save updated Seeker
// 	if err := h.DB.Save(seeker).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated personal info", "details": err.Error()})
// 		return
// 	}

// 	// Return updated personal info response
// 	c.JSON(http.StatusOK, dto.PersonalInfoResponse{
// 		AuthUserID:      seeker.AuthUserID,
// 		FirstName:       personalInfo.FirstName,
// 		SecondName:      personalInfo.SecondName,
// 		DateOfBirth:     personalInfo.DateOfBirth,
// 		Address:         personalInfo.Address,
// 		LinkedInProfile: personalInfo.LinkedInProfile,
// 	})
// }
