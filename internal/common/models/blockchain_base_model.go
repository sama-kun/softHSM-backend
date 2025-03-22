package models

import (
	"time"

	"github.com/google/uuid"
)

type BlockchainBaseModel struct {
	Id        uuid.UUID      `json:"id"`
	IsDeleted bool       `json:"isDeleted,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
}