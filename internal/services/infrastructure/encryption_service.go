package infrastructure

import (
	"gobi/pkg/utils"
)

// EncryptionServiceImpl implements EncryptionService
type EncryptionServiceImpl struct{}

// NewEncryptionService creates a new EncryptionService instance
func NewEncryptionService() *EncryptionServiceImpl {
	return &EncryptionServiceImpl{}
}

// Encrypt encrypts data using AES
func (s *EncryptionServiceImpl) Encrypt(data string) (string, error) {
	return utils.EncryptAES(data)
}

// Decrypt decrypts data using AES
func (s *EncryptionServiceImpl) Decrypt(data string) (string, error) {
	return utils.DecryptAES(data)
}
