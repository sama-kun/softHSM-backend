package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"soft-hsm/internal/blockchain-key/dto"
	"soft-hsm/internal/blockchain-key/models"
	"soft-hsm/internal/blockchain-key/repository"
	"soft-hsm/internal/blockchain-key/security"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"

	authservices "soft-hsm/internal/auth/services"
)

type BlockchainKeyServiceInterface interface {
	GenerateEthereumKey(ctx context.Context, userID int64, password string, key dto.GenerateKeyDTO) (*dto.GenerateKeyResponseDTO, error)
}

type BlockchainKeyService struct {
	blockchainKeyRepo repository.BlockchainKeyRepositoryInterface
	securityService   *security.SecurityService
	passwordService   *authservices.PasswordService
}

func NewBlockchainKeyService(blockchainKeyRepo repository.BlockchainKeyRepositoryInterface, securityService *security.SecurityService, passwordService *authservices.PasswordService) *BlockchainKeyService {
	return &BlockchainKeyService{
		blockchainKeyRepo: blockchainKeyRepo,
		securityService:   securityService,
		passwordService:   passwordService,
	}
}

func (s *BlockchainKeyService) GenerateEthereumKey(ctx context.Context, userID int64, password string, key dto.GenerateKeyDTO) (*dto.GenerateKeyResponseDTO, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации энтропии: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации мнемоники: %w", err)
	}

	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации master-ключа: %w", err)
	}

	ethKey, err := deriveEthereumKey(masterKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации Ethereum-ключа: %w", err)
	}

	publicKeyBytes := crypto.FromECDSAPub(&ethKey.PublicKey)
	address := crypto.PubkeyToAddress(ethKey.PublicKey).Hex()

	privateKeyBytes := crypto.FromECDSA(ethKey)
	encryptedKey, salt, err := s.securityService.EncryptPrivateKey(password, privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования приватного ключа: %w", err)
	}

	mnemonicHash := sha256.Sum256([]byte(mnemonic))

	blockchainKey := models.BlockchainKey{
		UserId:       userID,
		Name:         &key.Name,
		Description:  &key.Description,
		Blockchain:   models.Ethereum,
		Network:      "goerli",
		Address:      address,
		EncryptedKey: encryptedKey,
		PublicKey:    fmt.Sprintf("%x", publicKeyBytes),
		Salt:         salt,
		MnemonicHash: hex.EncodeToString(mnemonicHash[:]),
	}

	if _, err := s.blockchainKeyRepo.Save(ctx, &blockchainKey); err != nil {
		return nil, fmt.Errorf("ошибка сохранения ключа в БД: %w", err)
	}

	return &dto.GenerateKeyResponseDTO{
		Id:          blockchainKey.Id,
		Name:        blockchainKey.Name,
		Description: blockchainKey.Description,
		Blockchain:  blockchainKey.Blockchain,
		PublicKey:   blockchainKey.PublicKey,
		Address:     blockchainKey.Address,
		Mnemonic:    mnemonic,
		Network:     blockchainKey.Network,
	}, nil
}

func deriveEthereumKey(masterKey *hdkeychain.ExtendedKey) (*ecdsa.PrivateKey, error) {
	path := []uint32{
		hdkeychain.HardenedKeyStart + 44, // BIP-44
		hdkeychain.HardenedKeyStart + 60, // Ethereum (60)
		hdkeychain.HardenedKeyStart + 0,  // Аккаунт 0
		0,                                // Change 0 (обычно 0 для обычных адресов)
		0,                                // Индекс первого ключа
	}

	key := masterKey
	var err error
	for _, p := range path {
		key, err = key.Derive(p)
		if err != nil {
			return nil, fmt.Errorf("ошибка деривации ключа: %w", err)
		}
	}

	// Получаем приватный ключ
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения приватного ключа: %w", err)
	}

	return privateKey.ToECDSA(), nil
}

// ImportPrivateKey импортирует существующий приватный ключ или мнемонику
// func (b *BlockchainKeyService) ImportPrivateKey(ctx context.Context, userID int64, password, privateKeyHex, mnemonic string) (string, error) {
// 	var privateKey *ecdsa.PrivateKey
// 	var err error

// 	if mnemonic != "" {
// 		mnemonicHash := sha256.Sum256([]byte(mnemonic))
// 		existingKey, err := b.blockchainKeyRepo.FindByMnemonicHash(ctx, hex.EncodeToString(mnemonicHash[:]))
// 		if err != nil {
// 			return "", fmt.Errorf("мнемоника не найдена")
// 		}
// 		privateKey, err = crypto.HexToECDSA(existingKey.EncryptedKey)
// 	} else {
// 		privateKey, err = crypto.HexToECDSA(privateKeyHex)
// 	}

// 	if err != nil {
// 		return "", fmt.Errorf("неверный приватный ключ: %w", err)
// 	}

// 	// Получаем публичный ключ и Ethereum-адрес
// 	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
// 	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

// 	// Шифруем приватный ключ
// 	privateKeyBytes := crypto.FromECDSA(privateKey)
// 	encryptedKey, salt, err := b.securityService.EncryptPrivateKey(password, privateKeyBytes)
// 	if err != nil {
// 		return "", fmt.Errorf("ошибка шифрования приватного ключа: %w", err)
// 	}

// 	// Сохраняем в БД
// 	err = b.blockchainKeyRepo.Save(ctx, &repository.BlockchainKey{
// 		ID:           uuid.New().String(),
// 		UserID:       userID,
// 		Blockchain:   "ethereum",
// 		Address:      address,
// 		PublicKey:    fmt.Sprintf("%x", publicKeyBytes),
// 		EncryptedKey: encryptedKey,
// 		Salt:         salt,
// 	})

// 	if err != nil {
// 		return "", fmt.Errorf("ошибка сохранения в БД: %w", err)
// 	}

// 	return address, nil
// }
