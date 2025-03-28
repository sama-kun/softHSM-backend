package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"soft-hsm/internal/blockchain-key/dto"
	"soft-hsm/internal/blockchain-key/models"
	"soft-hsm/internal/blockchain-key/repository"
	"soft-hsm/internal/blockchain-key/security"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
)

// https://mainnet.infura.io/v3/6341f9c84e9c4f7a91cc518c69a5c11c
type BlockchainKeyServiceInterface interface {
	GenerateEthereumKey(ctx context.Context, userID int64, dto dto.GenerateKeyDTO) (*dto.GenerateKeyResponseDTO, error)
	ImportEthereumKey(ctx context.Context, userID int64, key dto.ImportKeyDTO) (*dto.ImportKeyResponseDTO, error)
	FindKeysByUserID(ctx context.Context, userID int64) ([]dto.SafeKeyResponseDTO, error)
	KeyDetail(ctx context.Context, keyID uuid.UUID, userID int64) (*dto.KeyDetailResponseDTO, error)
}

type BlockchainKeyService struct {
	blockchainKeyRepo repository.BlockchainKeyRepositoryInterface
	securityService   *security.SecurityService
	ethereumService   *EthereumService
}

func NewBlockchainKeyService(blockchainKeyRepo repository.BlockchainKeyRepositoryInterface, securityService *security.SecurityService, ethereumService *EthereumService) *BlockchainKeyService {
	return &BlockchainKeyService{
		blockchainKeyRepo: blockchainKeyRepo,
		securityService:   securityService,
		ethereumService:   ethereumService,
	}
}

// func (s *BlockchainKeyService) ImportEthereumKey(ctx context.Context, userID int64, key dto.ImportKeyDTO) (*dto.GenerateKeyResponseDTO, error) {
// 	// 1. Конвертация приватного ключа в ECDSA
// 	privateKey, err := crypto.HexToECDSA(key.PrivateKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка конвертации приватного ключа: %w", err)
// 	}

// 	// 2. Генерация публичного ключа
// 	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
// 	publicKeyHex := fmt.Sprintf("%x", publicKeyBytes)

// 	// 3. Генерация Ethereum-адреса
// 	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

// 	// 4. Шифрование приватного ключа
// 	privateKeyBytes := crypto.FromECDSA(privateKey)
// 	encryptedKey, salt, err := s.securityService.EncryptPrivateKey(privateKeyBytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка шифрования приватного ключа: %w", err)
// 	}

// 	// 5. Генерация хеша пустой мнемоники (если она не передается)
// 	mnemonicHash := sha256.Sum256([]byte(""))

// 	// 6. Сохранение в БД
// 	blockchainKey := models.BlockchainKey{
// 		UserId:       userID,
// 		Name:         &key.Name,
// 		Description:  &key.Description,
// 		Blockchain:   models.Ethereum,
// 		Network:      "goerli",
// 		Address:      address,
// 		EncryptedKey: encryptedKey,
// 		PublicKey:    publicKeyHex,
// 		Salt:         salt,
// 		MnemonicHash: hex.EncodeToString(mnemonicHash[:]), // Пустая мнемоника
// 	}

// 	if _, err := s.blockchainKeyRepo.Save(ctx, &blockchainKey); err != nil {
// 		return nil, fmt.Errorf("ошибка сохранения ключа в БД: %w", err)
// 	}

// 	// 7. Возвращаем результат
// 	return &dto.GenerateKeyResponseDTO{
// 		Id:          blockchainKey.Id,
// 		Name:        blockchainKey.Name,
// 		Description: blockchainKey.Description,
// 		Blockchain:  blockchainKey.Blockchain,
// 		PublicKey:   blockchainKey.PublicKey,
// 		Address:     blockchainKey.Address,
// 		Network:     blockchainKey.Network,
// 	}, nil
// }

func (s *BlockchainKeyService) KeyDetail(ctx context.Context, keyID uuid.UUID, userID int64) (*dto.KeyDetailResponseDTO, error) {
	key, err := s.blockchainKeyRepo.FindByID(ctx, keyID, userID)

	if err != nil {
		return nil, fmt.Errorf("cannot find key or not allowed")
	}

	fmt.Printf("Address: %s", key.Address)

	mainnetBalance, seplioBalance, err := s.ethereumService.GetBalances(ctx, key.Address)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &dto.KeyDetailResponseDTO{
		Id:             key.Id,
		Name:           key.Name,
		Description:    key.Description,
		Blockchain:     models.Ethereum,
		PublicKey:      key.PublicKey,
		Address:        key.Address,
		CreatedAt:      key.CreatedAt,
		MainnetBalance: mainnetBalance.Int64(),
		SeplioBalance:  seplioBalance.Int64(),
	}, nil
}

func (s *BlockchainKeyService) GenerateEthereumKey(ctx context.Context, userID int64, key dto.GenerateKeyDTO) (*dto.GenerateKeyResponseDTO, error) {
	// 1. Генерация энтропии и мнемоники
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации энтропии: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации мнемоники: %w", err)
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("ошибка: сгенерированная мнемоника недействительна")
	}
	fmt.Println("Мнемоника:", mnemonic)

	// 2. Генерация seed и master key
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации master-ключа: %w", err)
	}

	// 3. Деривация Ethereum-ключа
	ethKey, err := deriveEthereumKey(masterKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации Ethereum-ключа: %w", err)
	}

	// 4. Генерация публичного ключа
	if ethKey.PublicKey.X == nil || ethKey.PublicKey.Y == nil {
		return nil, fmt.Errorf("ошибка: публичный ключ Ethereum не инициализирован")
	}
	publicKeyBytes := crypto.FromECDSAPub(&ethKey.PublicKey)
	fmt.Println("Публичный ключ:", hex.EncodeToString(publicKeyBytes))

	// 5. Генерация адреса Ethereum
	address := crypto.PubkeyToAddress(ethKey.PublicKey).Hex()
	fmt.Println("Ethereum-адрес:", address)

	// 6. Шифрование приватного ключа
	privateKeyBytes := crypto.FromECDSA(ethKey)
	encryptedKey, salt, err := s.securityService.EncryptPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования приватного ключа: %w", err)
	}
	fmt.Println("Соль:", salt)

	// 7. Хеширование мнемоники
	mnemonicHash := sha256.Sum256([]byte(mnemonic))
	fmt.Println("Хеш мнемоники:", hex.EncodeToString(mnemonicHash[:]))

	// 8. Сохранение в БД
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

	// 9. Возвращаем результат
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

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения приватного ключа: %w", err)
	}

	return privateKey.ToECDSA(), nil
}

func (s *BlockchainKeyService) ImportEthereumKey(ctx context.Context, userID int64, key dto.ImportKeyDTO) (*dto.ImportKeyResponseDTO, error) {
	var privateKey *ecdsa.PrivateKey
	var mnemonicHash string

	// switch key.Type {
	// case dto.MnemonicKey:
	if !bip39.IsMnemonicValid(key.Input) {
		return nil, errors.New("недействительная мнемоника")
	}

	// Хешируем мнемонику
	hash := sha256.Sum256([]byte(key.Input))
	mnemonicHash = hex.EncodeToString(hash[:])

	// Проверяем в БД
	existingKey, err := s.blockchainKeyRepo.FindByMnemonicHash(ctx, mnemonicHash)
	if err == nil && existingKey != nil {
		return nil, errors.New("этот ключ уже импортирован")
	}

	// Генерируем ключ из мнемоники
	seed := bip39.NewSeed(key.Input, "")
	privateKey, err = crypto.ToECDSA(seed[:32])
	if err != nil {
		return nil, errors.New("ошибка генерации приватного ключа из мнемоники")
	}

	// case dto.PrivateKey:
	// 	// Проверяем валидность приватного ключа
	// 	privKeyBytes, err := hex.DecodeString(key.Input)
	// 	if err != nil || len(privKeyBytes) != 32 {
	// 		return nil, errors.New("недействительный приватный ключ")
	// 	}

	// 	privateKey, err = crypto.ToECDSA(privKeyBytes)
	// 	if err != nil {
	// 		return nil, errors.New("ошибка парсинга приватного ключа")
	// 	}

	// default:
	// 	return nil, errors.New("неверный тип ключа, используйте 'mnemonic' или 'key'")
	// }

	// Генерируем адрес
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	// Шифруем приватный ключ
	privateKeyBytes := crypto.FromECDSA(privateKey)
	encryptedKey, salt, err := s.securityService.EncryptPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, errors.New("ошибка шифрования приватного ключа")
	}

	// Сохраняем в БД
	blockchainKey := models.BlockchainKey{
		UserId:       userID,
		Blockchain:   models.Ethereum,
		Network:      "goerli",
		Address:      address,
		EncryptedKey: encryptedKey,
		PublicKey:    hex.EncodeToString(crypto.FromECDSAPub(publicKey)),
		Salt:         salt,
		MnemonicHash: mnemonicHash,
	}

	if _, err := s.blockchainKeyRepo.Save(ctx, &blockchainKey); err != nil {
		return nil, fmt.Errorf("ошибка сохранения ключа в БД: %w", err)
	}

	return &dto.ImportKeyResponseDTO{
		Address: address,
	}, nil
}

func (s *BlockchainKeyService) FindKeysByUserID(ctx context.Context, userID int64) ([]dto.SafeKeyResponseDTO, error) {
	keys, err := s.blockchainKeyRepo.FindByUserID(ctx, int64(userID))
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении ключей пользователя: %w", err)
	}

	var result []dto.SafeKeyResponseDTO
	for _, key := range keys {
		result = append(result, dto.SafeKeyResponseDTO{
			Id:          key.Id,
			Name:        key.Name,
			Description: key.Description,
			Blockchain:  key.Blockchain,
			PublicKey:   key.PublicKey,
			Address:     key.Address,
			Network:     key.Network,
		})
	}

	return result, nil
}
