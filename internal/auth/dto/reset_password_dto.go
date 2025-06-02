package dto

type ResetPasswordDTO struct {
	CurrentPassword    string `json:"currentPasswrod"`
	NewPassword        string `json:"newPasswrod"`
	ConfirmNewPassword string `json:"confirmNewPasswrod"`
}
