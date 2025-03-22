package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

type SecurityService struct {
	MasterKey []byte
}

func NewSecurityService() *SecurityService {
	masterPassword := os.Getenv("MASTER_PASSWORD")
	return &SecurityService{
		MasterKey: deriveMasterKey(masterPassword),
	}
}

func deriveMasterKey(password string) []byte {
	salt := []byte("fixed_salt_value_for_master_key")
	return pbkdf2.Key([]byte(password), salt, 500000, 32, sha256.New)
}

func (s *SecurityService) DeriveEncryptionKey(userPassword string, salt []byte) []byte {
	combinedKey := append(s.MasterKey, []byte(userPassword)...)
	return pbkdf2.Key(combinedKey, salt, 100000, 32, sha256.New)
}

func (s *SecurityService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func (s *SecurityService) EncryptPrivateKey(userPassword string, privateKey []byte) (string, string, error) {
	salt, err := s.GenerateSalt()
	if err != nil {
		return "", "", err
	}

	encryptionKey := s.DeriveEncryptionKey(userPassword, salt)

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	encrypted := aesGCM.Seal(nonce, nonce, privateKey, nil)
	return base64.StdEncoding.EncodeToString(encrypted), base64.StdEncoding.EncodeToString(salt), nil
}

func (s *SecurityService) DecryptPrivateKey(userPassword, encryptedKey, saltB64 string) ([]byte, error) {
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return nil, err
	}

	encryptionKey := s.DeriveEncryptionKey(userPassword, salt)

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	encryptedData, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("некорректные данные")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
