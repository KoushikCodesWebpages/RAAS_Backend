package features

import (
	"RAAS/models"
	"RAAS/dto"
	//"RAAS/middlewares"
	//"RAAS/security"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
)

// ProfileHandler struct
type ProfileHandler struct {
	DB *gorm.DB
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{DB: db}
}

// RetrieveProfile retrieves the profile of the authenticated user
func (h *ProfileHandler) RetrieveProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var seeker models.Seeker
	if err := h.DB.Preload("AuthUser").First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	profileResponse := dto.SeekerProfileResponse(seeker)
	c.JSON(http.StatusOK, profileResponse)
}

// UpdateProfile updates the profile of the authenticated user
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var updateData dto.SeekerResponse

	// Bind incoming JSON data
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var seeker models.Seeker
	if err := h.DB.Preload("AuthUser").First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Save updated profile
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update profile"})
		return
	}

	c.JSON(http.StatusOK, dto.SeekerProfileResponse(seeker))
}

// PatchProfile partially updates the profile of the authenticated user
func (h *ProfileHandler) PatchProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var updateData dto.SeekerResponse

	// Bind incoming JSON data
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var seeker models.Seeker
	if err := h.DB.Preload("AuthUser").First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	// Save updated profile
	if err := h.DB.Save(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not patch profile"})
		return
	}

	c.JSON(http.StatusOK, dto.SeekerProfileResponse(seeker))
}

// DeleteProfile deletes the profile of the authenticated user
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var seeker models.Seeker
	if err := h.DB.Preload("AuthUser").First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Delete the profile
	if err := h.DB.Delete(&seeker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}
