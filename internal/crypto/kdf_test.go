package crypto

import (
	"bytes"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt(16)
	if err != nil {
		t.Fatal(err)
	}
	if len(salt) != 16 {
		t.Fatalf("expected 16 bytes, got %d", len(salt))
	}

	salt2, _ := GenerateSalt(16)
	if bytes.Equal(salt, salt2) {
		t.Fatal("two salts should not be equal")
	}
}

func TestDeriveKey_Deterministic(t *testing.T) {
	password := []byte("test-password")
	salt := []byte("0123456789abcdef")
	p := DefaultKDFParams()

	key1 := DeriveKey(password, salt, p)
	key2 := DeriveKey(password, salt, p)

	if !bytes.Equal(key1, key2) {
		t.Fatal("same inputs should produce same key")
	}
	if len(key1) != int(p.KeyLen) {
		t.Fatalf("expected %d bytes, got %d", p.KeyLen, len(key1))
	}
}

func TestDeriveKey_DifferentPasswords(t *testing.T) {
	salt := []byte("0123456789abcdef")
	p := DefaultKDFParams()

	key1 := DeriveKey([]byte("password-a"), salt, p)
	key2 := DeriveKey([]byte("password-b"), salt, p)

	if bytes.Equal(key1, key2) {
		t.Fatal("different passwords should produce different keys")
	}
}
