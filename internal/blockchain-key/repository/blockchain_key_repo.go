package repository

import (
	"context"
	"fmt"
	"soft-hsm/internal/blockchain-key/models"
	"soft-hsm/internal/storage"

	"github.com/google/uuid"
)

type BlockchainKeyRepositoryInterface interface {
	Save(ctx context.Context, blockchainKey *models.BlockchainKey) (*models.BlockchainKey, error)
	ImportKey(ctx context.Context, key *models.BlockchainKey) (*models.BlockchainKey, error)
	FindByIDWithKey(ctx context.Context, id uuid.UUID) (*models.BlockchainKey, error)
	FindByMnemonicHash(ctx context.Context, mnemonicHash string) (*models.BlockchainKey, error)
	FindByUserID(ctx context.Context, userID int64) ([]models.BlockchainKey, error)
	FindByID(ctx context.Context, id uuid.UUID, userID int64) (*models.BlockchainKey, error)
	DeleteEthereumKeyByID(ctx context.Context, id uuid.UUID, userID int64) error
}

type BlockchainKeyRepository struct {
	db *storage.Postgres
}

func NewBlockchainKeyRepository(db *storage.Postgres) BlockchainKeyRepositoryInterface {
	return &BlockchainKeyRepository{db: db}
}

func (r *BlockchainKeyRepository) DeleteEthereumKeyByID(ctx context.Context, id uuid.UUID, userID int64) error {
	query := `DELETE FROM blockchain_keys WHERE id = $1 AND user_id = $2`
	_, err := r.db.Conn().Exec(ctx, query, id, userID)
	return err
}

func (r *BlockchainKeyRepository) FindByID(ctx context.Context, id uuid.UUID, userID int64) (*models.BlockchainKey, error) {
	const fn = "blockchain_key_repo.FindByID"

	var key models.BlockchainKey

	query := `
	  SELECT id, blockchain, user_id, name, description, network, address, public_key
		FROM blockchain_keys
		WHERE id = $1 AND user_id = $2
	`

	err := r.db.Conn().QueryRow(ctx, query, id, userID).Scan(
		&key.Id,
		&key.Blockchain,
		&key.UserId,
		&key.Name,
		&key.Description,
		&key.Network,
		&key.Address,
		&key.PublicKey,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain key /%s/: %w", fn, err)
	}

	return &key, nil
}

func (r *BlockchainKeyRepository) FindByIDWithKey(ctx context.Context, id uuid.UUID) (*models.BlockchainKey, error) {
	const fn = "blockchain_key_repo.FindByIDWithKey"

	var key models.BlockchainKey

	query := `
	  SELECT id, blockchain, user_id, encrypted_key, name, description, network, address, salt
		FROM blockchain_keys
		WHERE id = $1
	`

	err := r.db.Conn().QueryRow(ctx, query, id).Scan(
		&key.Id,
		&key.Blockchain,
		&key.UserId,
		&key.EncryptedKey,
		&key.Name,
		&key.Description,
		&key.Network,
		&key.Address,
		&key.Salt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain key /%s/: %w", fn, err)
	}

	return &key, nil
}

func (r *BlockchainKeyRepository) Save(ctx context.Context, blockchainKey *models.BlockchainKey) (*models.BlockchainKey, error) {
	const fn = "blockchain_key_repo.Save"
	tx, err := r.db.Conn().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction /%s/: %w", fn, err)
	}
	defer tx.Rollback(ctx)

	query := `
	INSERT INTO blockchain_keys 
	(user_id, blockchain, address, encrypted_key, public_key, name, description, network, salt, mnemonic_hash) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, name, description, address, encrypted_key, public_key, network, salt, mnemonic_hash
`

	err = tx.QueryRow(ctx, query,
		blockchainKey.UserId,
		blockchainKey.Blockchain,
		blockchainKey.Address,
		blockchainKey.EncryptedKey,
		blockchainKey.PublicKey,
		blockchainKey.Name,
		blockchainKey.Description,
		"goerli",
		blockchainKey.Salt,
		blockchainKey.MnemonicHash,
	).Scan(&blockchainKey.Id, &blockchainKey.Name, &blockchainKey.Description,
		&blockchainKey.Address, &blockchainKey.EncryptedKey, &blockchainKey.PublicKey,
		&blockchainKey.Network, &blockchainKey.Salt, &blockchainKey.MnemonicHash)

	if err != nil {
		return nil, fmt.Errorf("failed to insert blockchain key: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return blockchainKey, nil
}

func (r *BlockchainKeyRepository) ImportKey(ctx context.Context, key *models.BlockchainKey) (*models.BlockchainKey, error) {
	return r.Save(ctx, key)
}

func (r *BlockchainKeyRepository) GetPublicKey(ctx context.Context, keyId uuid.UUID) (string, error) {
	const fn = "blockchain_key_repo.GetPublicKey"

	query := `SELECT address FROM blockchain_keys WHERE id = $1 AND is_deleted = FALSE`
	var address string
	err := r.db.Conn().QueryRow(ctx, query, keyId).Scan(&address)
	if err != nil {
		return "", fmt.Errorf("failed to get public key /%s/: %w", fn, err)
	}
	return address, nil
}

func (r *BlockchainKeyRepository) FindByMnemonicHash(ctx context.Context, mnemonicHash string) (*models.BlockchainKey, error) {
	const fn = "blockchain_key_repo.FindByMnemonicHash"

	var key models.BlockchainKey

	query := `
	  SELECT id, blockchain, user_id, encrypted_key, name, description, network, address, mnemonic_hash
		FROM blockchain_keys
		WHERE mnemonic_hash = $1
	`

	err := r.db.Conn().QueryRow(ctx, query, mnemonicHash).Scan(
		&key.Id,
		&key.Blockchain,
		&key.UserId,
		&key.EncryptedKey,
		&key.Name,
		&key.Description,
		&key.Network,
		&key.Address,
		&key.MnemonicHash,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find blockchain key by mnemonic hash /%s/: %w", fn, err)
	}

	return &key, nil
}

func (r *BlockchainKeyRepository) FindByUserID(ctx context.Context, userID int64) ([]models.BlockchainKey, error) {
	const fn = "blockchain_key_repo.FindByUserID"

	query := `
		SELECT id, blockchain, user_id, name, description, network, address, public_key
		FROM blockchain_keys
		WHERE user_id = $1
	`

	rows, err := r.db.Conn().Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain keys /%s/: %w", fn, err)
	}
	defer rows.Close()

	var keys []models.BlockchainKey
	for rows.Next() {
		var key models.BlockchainKey
		if err := rows.Scan(
			&key.Id,
			&key.Blockchain,
			&key.UserId,
			&key.Name,
			&key.Description,
			&key.Network,
			&key.Address,
			&key.PublicKey,
		); err != nil {
			return nil, fmt.Errorf("failed to scan blockchain key /%s/: %w", fn, err)
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over blockchain keys /%s/: %w", fn, err)
	}

	return keys, nil
}
