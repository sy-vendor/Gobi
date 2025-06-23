package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gobi/config"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT generates a new JWT token for a user
func GenerateJWT(userID uint, role string) (string, error) {
	cfg := config.AppConfig
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(cfg.JWT.ExpirationHours) * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// EncryptAES 加密明文，返回 base64 字符串
func EncryptAES(plaintext string) (string, error) {
	key := []byte(os.Getenv("DATA_SOURCE_SECRET"))
	if len(key) == 0 {
		return "", fmt.Errorf("DATA_SOURCE_SECRET environment variable not set")
	}
	if len(key) != 32 {
		return "", fmt.Errorf("DATA_SOURCE_SECRET must be 32 bytes (256 bit), but got %d bytes", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES 解密 base64 字符串，返回明文
func DecryptAES(ciphertext string) (string, error) {
	key := []byte(os.Getenv("DATA_SOURCE_SECRET"))
	if len(key) == 0 {
		return "", fmt.Errorf("DATA_SOURCE_SECRET environment variable not set")
	}
	if len(key) != 32 {
		return "", fmt.Errorf("DATA_SOURCE_SECRET must be 32 bytes (256 bit), but got %d bytes", len(key))
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertextBytes := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
