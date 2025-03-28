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

func NewSecurityService() (*SecurityService, error) {
	masterKey, err := loadMasterKeyFromFile("master_key.enc")
	if err != nil {
		return nil, err
	}
	return &SecurityService{
		MasterKey: masterKey,
	}, nil
}

func loadMasterKeyFromFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *SecurityService) DeriveEncryptionKey(salt []byte) []byte {
	return pbkdf2.Key(s.MasterKey, salt, 100000, 32, sha256.New)
}

func (s *SecurityService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func (s *SecurityService) EncryptPrivateKey(privateKey []byte) (string, string, error) {
	salt, err := s.GenerateSalt()
	if err != nil {
		return "", "", err
	}
	encryptionKey := s.DeriveEncryptionKey(salt)

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

func (s *SecurityService) DecryptPrivateKey(encryptedKey, saltB64 string) ([]byte, error) {
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return nil, err
	}
	encryptionKey := s.DeriveEncryptionKey(salt)

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
