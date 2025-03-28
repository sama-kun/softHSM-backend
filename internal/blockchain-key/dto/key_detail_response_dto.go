package dto

import (
	"soft-hsm/internal/blockchain-key/models"
	"time"

	"github.com/google/uuid"
)

type KeyDetailResponseDTO struct {
	Id             uuid.UUID             `json:"id"`
	Name           *string               `json:"name,omitempty"`
	Description    *string               `json:"description,omitempty"`
	Blockchain     models.BlockchainType `json:"blockchain"`
	PublicKey      string                `json:"publicKey"`
	Address        string                `json:"address"`
	CreatedAt      time.Time             `json:"createdAt,omitempty"`
	MainnetBalance int64                 `json:"mainnetBalance"`
	SeplioBalance  int64                 `json:"seplioBalance"`
}
