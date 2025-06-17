package dto

import "soft-hsm/internal/user/models"

type LoginResponseDTO struct {
	SessionToken string       `json:"sessionToken"`
	User         *models.User `json:"user"`
}
