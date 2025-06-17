package dto

type ImportPrivateKeyDTO struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	PrivateKeyBase64 string `json:"privateKey"`
}
