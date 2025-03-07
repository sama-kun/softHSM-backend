package dto

type ActivateDTO struct {
	ActivateToken string `json:"activateToken" validate:"required"`
}
