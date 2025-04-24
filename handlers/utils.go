package handlers

import (
	"RAAS/models"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"encoding/json"
)

// FindSeekerByUserID is a global utility function to find a Seeker by userID
func FindSeekerByUserID(db *gorm.DB, userID uuid.UUID) (*models.Seeker, error) {
	var seeker models.Seeker
	if err := db.First(&seeker, "auth_user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &seeker, nil
}


func IsFieldFilled(field []byte) (bool, error) {
	if len(field) > 0 {
		var existingField map[string]interface{}
		if err := json.Unmarshal(field, &existingField); err == nil && len(existingField) > 0 {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}