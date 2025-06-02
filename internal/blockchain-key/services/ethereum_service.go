package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	infuraMainnet = "https://mainnet.infura.io/v3/6341f9c84e9c4f7a91cc518c69a5c11c"
	infuraSepolia = "https://sepolia.infura.io/v3/6341f9c84e9c4f7a91cc518c69a5c11c"
)

type EthereumService struct{}

func NewEthereumService() *EthereumService {

	return &EthereumService{}
}

func (s *EthereumService) GetBalance(ctx context.Context, infuraURL, address string) (*big.Int, error) {
	client, err := rpc.DialContext(ctx, infuraURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var balanceHex string
	err = client.Call(&balanceHex, "eth_getBalance", common.HexToAddress(address), "latest")
	if err != nil {
		return nil, err
	}

	fmt.Println(balanceHex)

	balance := new(big.Int)
	balance.SetString(balanceHex[2:], 16) // Конвертация из hex в big.Int
	return balance, nil
}

func (s *EthereumService) GetBalances(ctx context.Context, address string) (*big.Int, *big.Int, error) {
	mainnetBalance, err := s.GetBalance(ctx, infuraMainnet, address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get mainnet balance: %w", err)
	}

	sepoliaBalance, err := s.GetBalance(ctx, infuraSepolia, address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sepolia balance: %w", err)
	}

	return mainnetBalance, sepoliaBalance, nil
}
