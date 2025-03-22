package models

import baseModels "soft-hsm/internal/common/models"

type BlockchainKey struct {
	baseModels.BlockchainBaseModel
	Name         *string        `json:"name,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Network      string         `json:"network"`
	UserId       int64          `json:"userId"`
	Blockchain   BlockchainType `json:"blockchain"`
	Address      string         `json:"address"`
	EncryptedKey string         `json:"-"`
	PublicKey    string         `json:"publicKey"`
	Salt         string         `json:"-"`

	MnemonicHash string `json:"-"`
}
