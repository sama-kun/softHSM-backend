package dto

type SetMasterPassword struct {
	SessionToken string `json:"sessionToken"`
	Otp          string `json:"otp"`
}
