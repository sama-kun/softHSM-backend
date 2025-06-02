package dto

import (
	"errors"
	"math/big"

	"github.com/google/uuid"
)

type SendEthereumDTO struct {
	KeyId       uuid.UUID `json:"keyId"`
	ToAddress   string    `json:"toAddress"`
	AmountInWei string    `json:"amountInWei"`
}

func (dto *SendEthereumDTO) ToBigInt() (*big.Int, error) {
	amount := new(big.Int)
	_, ok := amount.SetString(dto.AmountInWei, 10)
	if !ok {
		return nil, errors.New("invalid big.Int format")
	}
	return amount, nil
}
