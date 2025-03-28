package dto

type KeyType string

const (
	MnemonicKey KeyType = "mnemonic"
	PrivateKey  KeyType = "key"
)

type ImportKeyDTO struct {
	Type  KeyType `json:"type" binding:"required,oneof=mnemonic key"`
	Input string  `json:"input" binding:"required"`
}

type ImportKeyResponseDTO struct {
	Address string `json:"address"`
}
