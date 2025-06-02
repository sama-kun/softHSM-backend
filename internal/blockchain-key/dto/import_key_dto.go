package dto

type KeyType string

const (
	MnemonicKey KeyType = "mnemonic"
	PrivateKey  KeyType = "key"
)

type ImportKeyDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PrivateKey  string `json:"privateKey"`
}

type ImportKeyResponseDTO struct {
	Address string `json:"address"`
}
