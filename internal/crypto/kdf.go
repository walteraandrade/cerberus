package crypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type KDFParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLen     int
	KeyLen      uint32
}

func DefaultKDFParams() KDFParams {
	return KDFParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLen:     16,
		KeyLen:      32,
	}
}

func GenerateSalt(n int) ([]byte, error) {
	salt := make([]byte, n)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}
	return salt, nil
}

func DeriveKey(password, salt []byte, p KDFParams) []byte {
	return argon2.IDKey(password, salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLen)
}
