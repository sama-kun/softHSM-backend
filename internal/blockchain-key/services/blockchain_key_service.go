package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"soft-hsm/internal/blockchain-key/dto"
	"soft-hsm/internal/blockchain-key/models"
	"soft-hsm/internal/blockchain-key/repository"
	"soft-hsm/internal/blockchain-key/security"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
)

// https://mainnet.infura.io/v3/6341f9c84e9c4f7a91cc518c69a5c11c
type BlockchainKeyServiceInterface interface {
	GenerateEthereumKey(ctx context.Context, userID int64, dto dto.GenerateKeyDTO) (*dto.GenerateKeyResponseDTO, error)
	ImportEthereumKey(ctx context.Context, userID int64, key dto.ImportKeyDTO) (*dto.GenerateKeyResponseDTO, error)
	FindKeysByUserID(ctx context.Context, userID int64) ([]dto.SafeKeyResponseDTO, error)
	KeyDetail(ctx context.Context, keyID uuid.UUID, userID int64) (*dto.KeyDetailResponseDTO, error)
	SendEthereumTransaction(
		ctx context.Context,
		userID int64,
		keyID uuid.UUID,
		toAddress string,
		amountInWei *big.Int,
	) (string, error)
	ExportAndDeleteEthereumKeyByID(ctx context.Context, id uuid.UUID, userID int64) ([]byte, error)
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
	// fmt.Println("Соль:", salt)

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

func (s *BlockchainKeyService) SendEthereumTransaction(
	ctx context.Context,
	userID int64,
	keyID uuid.UUID,
	toAddress string,
	amountInWei *big.Int,
) (string, error) {
	// 1. Получаем ключ из БД
	key, err := s.blockchainKeyRepo.FindByIDWithKey(ctx, keyID)
	if err != nil {
		return "", fmt.Errorf("ключ не найден или доступ запрещен: %w", err)
	}

	// 2. Расшифровываем приватный ключ
	privateKeyBytes, err := s.securityService.DecryptPrivateKey(key.EncryptedKey, key.Salt)
	if err != nil {
		return "", fmt.Errorf("ошибка расшифровки приватного ключа: %w", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования ключа: %w", err)
	}

	// 3. Подключение к публичному RPC Ethereum Mainnet
	client, err := ethclient.Dial("https://rpc.ankr.com/eth/c53f01bd6d97a83e7e6cdb10ac9b4b3438186f73089a0e5d4d495bec5b1a616b")
	if err != nil {
		return "", fmt.Errorf("ошибка подключения к сети Ethereum: %w", err)
	}
	defer client.Close()

	// 4. Получаем адрес отправителя
	fromAddress := common.HexToAddress(key.Address)
	if !common.IsHexAddress(toAddress) {
		return "", fmt.Errorf("некорректный адрес получателя")
	}
	to := common.HexToAddress(toAddress)

	// 5. Получаем nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("ошибка получения nonce: %w", err)
	}

	// 6. Получаем цену газа и увеличиваем для приоритета
	suggestedGasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("ошибка получения цены газа: %w", err)
	}
	// Увеличим gas price в 2 раза для приоритета
	// gasPrice := new(big.Int).Mul(suggestedGasPrice, big.NewInt(2))

	// 7. Собираем транзакцию
	gasLimit := uint64(21000)

	tx := types.NewTransaction(nonce, to, amountInWei, gasLimit, suggestedGasPrice, nil)

	// 8. Подписываем транзакцию для Mainnet
	chainID := big.NewInt(1)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("ошибка подписи транзакции: %w", err)
	}

	// 9. Отправляем транзакцию в публичный mempool
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки транзакции: %w", err)
	}

	// 10. Лог транзакции
	fmt.Println("📤 Ethereum Mainnet Transaction Sent:")
	fmt.Printf("🔑 From: %s\n", fromAddress.Hex())
	fmt.Printf("📥 To:   %s\n", to.Hex())
	fmt.Printf("💰 Amount: %s Wei\n", amountInWei.String())
	// fmt.Printf("⛽ GasPrice: %s Wei\n", gasPrice.String())
	fmt.Printf("🔢 Nonce: %d\n", nonce)
	fmt.Printf("🔗 TxHash: %s\n", signedTx.Hash().Hex())

	// 11. Возвращаем хэш транзакции
	return signedTx.Hash().Hex(), nil
}

func (s *BlockchainKeyService) ImportEthereumKey(ctx context.Context, userID int64, key dto.ImportKeyDTO) (*dto.GenerateKeyResponseDTO, error) {
	// 1. Проверка валидности мнемоники
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации энтропии: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации мнемоники: %w", err)
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("мнемоника недействительна")
	}

	// 2. Генерация сидов и master key
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания master-ключа: %w", err)
	}

	// 3. Деривация Ethereum ключа
	ethKey, err := deriveEthereumKey(masterKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка деривации Ethereum-ключа: %w", err)
	}

	// 4. Проверка публичного ключа
	if ethKey.PublicKey.X == nil || ethKey.PublicKey.Y == nil {
		return nil, fmt.Errorf("ошибка: публичный ключ Ethereum не инициализирован")
	}

	publicKeyBytes := crypto.FromECDSAPub(&ethKey.PublicKey)
	address := crypto.PubkeyToAddress(ethKey.PublicKey).Hex()

	// 5. Шифрование приватного ключа
	privateKeyBytes := crypto.FromECDSA(ethKey)
	encryptedKey, salt, err := s.securityService.EncryptPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования приватного ключа: %w", err)
	}

	// 6. Хеш мнемоники
	mnemonicHash := sha256.Sum256([]byte(mnemonic))

	// 7. Сохранение в БД
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
		return nil, fmt.Errorf("ошибка при сохранении импортированного ключа: %w", err)
	}

	// 8. Ответ
	return &dto.GenerateKeyResponseDTO{
		Id:          blockchainKey.Id,
		Name:        blockchainKey.Name,
		Description: blockchainKey.Description,
		Blockchain:  blockchainKey.Blockchain,
		Address:     blockchainKey.Address,
		PublicKey:   blockchainKey.PublicKey,
		Mnemonic:    mnemonic,
		Network:     blockchainKey.Network,
	}, nil
}

func (s *BlockchainKeyService) ExportAndDeleteEthereumKeyByID(ctx context.Context, id uuid.UUID, userID int64) ([]byte, error) {
	// 1. Получаем ключ из базы
	key, err := s.blockchainKeyRepo.FindByIDWithKey(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	// 2. Дешифруем приватный ключ
	decryptedKey, err := s.securityService.DecryptPrivateKey(key.EncryptedKey, key.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}

	// 3. Удаляем ключ из базы
	err = s.blockchainKeyRepo.DeleteEthereumKeyByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete key: %w", err)
	}

	return decryptedKey, nil
}
