package models

import userModels "soft-hsm/internal/user/models"

type BlockchainKeyWithUser struct {
	BlockchainKey
	User userModels.User
}