package storage

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/walteraandrade/cerberus/internal/crypto"
)

const (
	headerMagic   = "CERB"
	headerVersion = 1
)

type VaultFile struct {
	Salt       []byte
	WrappedKey []byte
	Ciphertext []byte
}

func Write(path string, salt, wrappedKey, ciphertext []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".vault-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpPath)
	}()

	if _, err := tmp.Write([]byte(headerMagic)); err != nil {
		return err
	}
	if err := binary.Write(tmp, binary.BigEndian, uint8(headerVersion)); err != nil {
		return err
	}
	if err := writeChunk(tmp, salt); err != nil {
		return err
	}
	if err := writeChunk(tmp, wrappedKey); err != nil {
		return err
	}
	if err := writeChunk(tmp, ciphertext); err != nil {
		return err
	}

	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, 0600); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func Read(path string) (*VaultFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) < 5 {
		return nil, fmt.Errorf("vault file too small")
	}
	if string(data[:4]) != headerMagic {
		return nil, fmt.Errorf("invalid vault file magic")
	}
	if data[4] != headerVersion {
		return nil, fmt.Errorf("unsupported vault version %d", data[4])
	}

	pos := 5
	salt, n, err := readChunk(data, pos)
	if err != nil {
		return nil, fmt.Errorf("read salt: %w", err)
	}
	pos += n

	wrappedKey, n, err := readChunk(data, pos)
	if err != nil {
		return nil, fmt.Errorf("read wrapped key: %w", err)
	}
	pos += n

	ciphertext, _, err := readChunk(data, pos)
	if err != nil {
		return nil, fmt.Errorf("read ciphertext: %w", err)
	}

	return &VaultFile{
		Salt:       salt,
		WrappedKey: wrappedKey,
		Ciphertext: ciphertext,
	}, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CreateVault(path string, password []byte, plaintext []byte, params crypto.KDFParams) error {
	salt, err := crypto.GenerateSalt(params.SaltLen)
	if err != nil {
		return err
	}

	wrappingKey := crypto.DeriveKey(password, salt, params)
	defer crypto.ZeroBytes(wrappingKey)

	vaultKey, err := crypto.GenerateVaultKey(int(params.KeyLen))
	if err != nil {
		return err
	}

	wrappedKey, err := crypto.WrapKey(wrappingKey, vaultKey)
	if err != nil {
		crypto.ZeroBytes(vaultKey)
		return err
	}

	ciphertext, err := crypto.Encrypt(vaultKey, plaintext)
	crypto.ZeroBytes(vaultKey)
	if err != nil {
		return err
	}

	return Write(path, salt, wrappedKey, ciphertext)
}

func OpenVault(path string, password []byte, params crypto.KDFParams) ([]byte, error) {
	vf, err := Read(path)
	if err != nil {
		return nil, err
	}

	wrappingKey := crypto.DeriveKey(password, vf.Salt, params)
	defer crypto.ZeroBytes(wrappingKey)

	vaultKey, err := crypto.UnwrapKey(wrappingKey, vf.WrappedKey)
	if err != nil {
		return nil, fmt.Errorf("wrong password or corrupted vault")
	}

	plaintext, err := crypto.Decrypt(vaultKey, vf.Ciphertext)
	crypto.ZeroBytes(vaultKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt vault: %w", err)
	}

	return plaintext, nil
}

func writeChunk(f *os.File, data []byte) error {
	if err := binary.Write(f, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err := f.Write(data)
	return err
}

func readChunk(data []byte, pos int) ([]byte, int, error) {
	if pos+4 > len(data) {
		return nil, 0, fmt.Errorf("unexpected EOF reading chunk length")
	}
	size := int(binary.BigEndian.Uint32(data[pos : pos+4]))
	pos += 4
	if pos+size > len(data) {
		return nil, 0, fmt.Errorf("unexpected EOF reading chunk data")
	}
	return data[pos : pos+size], 4 + size, nil
}
