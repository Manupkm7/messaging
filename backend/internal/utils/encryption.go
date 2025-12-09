package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"backend/internal/config"
)

// GenerateAESKey genera una clave AES-256 desde el JWT secret
func GenerateAESKey() ([]byte, error) {
	cfg := config.AppConfig
	hash := sha256.Sum256([]byte(cfg.JWTSecret))
	return hash[:], nil
}

// EncryptMessage cifra un mensaje usando AES-256-GCM
func EncryptMessage(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key, err := GenerateAESKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptMessage descifra un mensaje cifrado con AES-256-GCM
func DecryptMessage(encrypted string) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	key, err := GenerateAESKey()
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", errors.New("invalid encrypted message format")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("failed to decrypt message")
	}

	return string(plaintext), nil
}

// EncryptData cifra datos binarios usando AES-256-GCM
func EncryptData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	key, err := GenerateAESKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// DecryptData descifra datos binarios cifrados con AES-256-GCM
func DecryptData(encrypted []byte) ([]byte, error) {
	if len(encrypted) == 0 {
		return encrypted, nil
	}

	key, err := GenerateAESKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, errors.New("encrypted data too short")
	}

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt data")
	}

	return plaintext, nil
}
