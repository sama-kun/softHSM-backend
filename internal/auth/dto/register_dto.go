package dto

type RegisterDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Login    string `json:"login" validate:"required,min=5"`
}
