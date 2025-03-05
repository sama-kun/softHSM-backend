package models

import "soft-hsm/internal/common/models"

type User struct {
	models.BaseModel
	Email    string `json:"email"`
	Password string `json:"-"`
	IsActive bool   `json:"is_active"`
	Login    string `json:"login"`
}
