package dto

type CheckMasterPasswordResponseDTO struct {
	SessionToken string `json:"sessionToken"`
	Id           int64  `json:"id"`
}
