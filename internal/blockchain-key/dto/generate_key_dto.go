package dto

import "soft-hsm/internal/blockchain-key/models"

type GenerateKeyDTO struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Blockchain  models.BlockchainType `json:"blockchain"`
}
