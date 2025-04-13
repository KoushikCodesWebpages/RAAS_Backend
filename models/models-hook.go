package models

import (
	"gorm.io/gorm"
	"fmt"
	"github.com/google/uuid"
)

const (
    RoleAdmin = "admin"
    RoleUser  = "seeker"
)

//Auth User
// Before saving, validate the role
func (u *AuthUser) BeforeCreate(tx *gorm.DB) (err error) {
    if u.Role != RoleAdmin && u.Role != RoleUser {
        return fmt.Errorf("role must be either 'admin' or 'user'")
    }
    return nil
}

func (s *Seeker) BeforeSave(tx *gorm.DB) (err error) {
	if s.SubscriptionTier != "free" && s.SubscriptionTier != "premium" {
		return fmt.Errorf("invalid subscription tier: %s", s.SubscriptionTier)
	}
	return nil
}

func (timeline *UserEntryTimeline) BeforeCreate(tx *gorm.DB) (err error) {
    timeline.ID = uuid.New()
    return
}
