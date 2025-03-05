package dto

import "soft-hsm/internal/user/models"

type LoginResponseDTO struct {
	AccessToken string       `json:"accessToken"`
	User        *models.User `json:"user"`
}
