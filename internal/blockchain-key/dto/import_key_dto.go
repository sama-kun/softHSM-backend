package dto

type KeyType string

const (
	MnemonicKey KeyType = "mnemonic"
	PrivateKey  KeyType = "key"
)

type RecoverKeyDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Mnemonic    string `json:"mnemonic"`
}

type ImportKeyResponseDTO struct {
	Address string `json:"address"`
}
