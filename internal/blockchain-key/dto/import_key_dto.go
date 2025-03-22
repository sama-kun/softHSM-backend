package dto

import "soft-hsm/internal/blockchain-key/models"

type ImportKeyDTO struct {
	PrivateKey  string                `json:"privateKey"`
	Blockchain  models.BlockchainType `json:"blockchain"`
	Network     string                `json:"network"`
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
}
