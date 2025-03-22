package models


type BlockchainType string

const (
	Ethereum BlockchainType = "ethereum"
	Bitcoin  BlockchainType = "bitcoin"
	Solana   BlockchainType = "solana"
)

func IsValidBlockchain(b BlockchainType) bool {
	switch b {
	case Ethereum, Bitcoin, Solana:
		return true
	}
	return false
}