package models

import (
	"gorm.io/gorm"
	"fmt"
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
