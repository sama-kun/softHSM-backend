package dto

import (
	"soft-hsm/internal/blockchain-key/models"

	"github.com/google/uuid"
)

type SafeKeyResponseDTO struct {
	Id          uuid.UUID             `json:"id"`
	Name        *string               `json:"name,omitempty"`
	Description *string               `json:"description,omitempty"`
	Blockchain  models.BlockchainType `json:"blockchain"`
	PublicKey   string                `json:"publicKey"`
	Address     string                `json:"address"`
	Network     string                `json:"network"`
}
