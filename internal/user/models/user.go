package models

import "soft-hsm/internal/common/models"

type User struct {
	models.BaseModel
	Email          string `json:"email"`
	Password       string `json:"-"`
	Login          string `json:"login"`
	MasterPassword string `json:"-"`

	IsActive       bool `json:"is_active"`
	IsVerified     bool `json:"is_verified" db:"is_verified" default:"false"`
	IsActiveMaster bool `json:"is_active_master" db:"is_active_master" default:"false"`
	// IsActiveFaceID bool `json:"is_active_faceid" db:"is_active_faceid" default:"false"`
}
