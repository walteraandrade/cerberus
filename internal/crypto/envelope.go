package crypto

import (
	"crypto/rand"
	"fmt"
)

func GenerateVaultKey(keyLen int) ([]byte, error) {
	key := make([]byte, keyLen)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("generate vault key: %w", err)
	}
	return key, nil
}

func WrapKey(wrappingKey, vaultKey []byte) ([]byte, error) {
	return Encrypt(wrappingKey, vaultKey)
}

func UnwrapKey(wrappingKey, wrappedKey []byte) ([]byte, error) {
	return Decrypt(wrappingKey, wrappedKey)
}
