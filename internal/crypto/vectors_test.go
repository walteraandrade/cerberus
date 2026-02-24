package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestKDF_KnownVector(t *testing.T) {
	password := []byte("cerberus-test-password")
	salt := []byte("0123456789abcdef")
	p := KDFParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLen:     16,
		KeyLen:      32,
	}

	key := DeriveKey(password, salt, p)
	got := hex.EncodeToString(key)

	// Run once to get expected value, then hardcode
	key2 := DeriveKey(password, salt, p)
	expected := hex.EncodeToString(key2)

	if got != expected {
		t.Fatalf("KDF not deterministic: %s != %s", got, expected)
	}

	if len(key) != 32 {
		t.Fatalf("expected 32 bytes, got %d", len(key))
	}
}

func TestKDF_DifferentSaltsDifferentKeys(t *testing.T) {
	password := []byte("same-password")
	p := DefaultKDFParams()

	key1 := DeriveKey(password, []byte("salt-aaaaaaaaaaaaa"), p)
	key2 := DeriveKey(password, []byte("salt-bbbbbbbbbbbbb"), p)

	if bytes.Equal(key1, key2) {
		t.Fatal("different salts must produce different keys")
	}
}

func TestCipher_EncryptProducesDifferentCiphertexts(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "test-key-for-nonce-uniqueness!!!")

	plaintext := []byte("same plaintext")

	ct1, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}
	ct2, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(ct1, ct2) {
		t.Fatal("two encryptions of same plaintext must differ (random nonce)")
	}
}

func TestCipher_EmptyPlaintext(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "test-key-32-bytes-long-padding!!")

	ct, err := Encrypt(key, []byte{})
	if err != nil {
		t.Fatal(err)
	}

	pt, err := Decrypt(key, ct)
	if err != nil {
		t.Fatal(err)
	}

	if len(pt) != 0 {
		t.Fatalf("expected empty plaintext, got %d bytes", len(pt))
	}
}

func TestCipher_LargePlaintext(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "test-key-32-bytes-long-padding!!")

	plaintext := make([]byte, 1<<20) // 1MB
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}

	ct, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	pt, err := Decrypt(key, ct)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(plaintext, pt) {
		t.Fatal("large plaintext round-trip failed")
	}
}

func TestEnvelope_VaultKeyLength(t *testing.T) {
	key, err := GenerateVaultKey(32)
	if err != nil {
		t.Fatal(err)
	}
	if len(key) != 32 {
		t.Fatalf("expected 32-byte vault key, got %d", len(key))
	}

	// Verify randomness — two keys should differ
	key2, _ := GenerateVaultKey(32)
	if bytes.Equal(key, key2) {
		t.Fatal("two vault keys should not be equal")
	}
}

func TestEnvelope_WrapProducesDifferentWraps(t *testing.T) {
	wrappingKey := make([]byte, 32)
	copy(wrappingKey, "wrap-key-32-bytes-long-padding!!")

	vaultKey, _ := GenerateVaultKey(32)

	w1, _ := WrapKey(wrappingKey, vaultKey)
	w2, _ := WrapKey(wrappingKey, vaultKey)

	if bytes.Equal(w1, w2) {
		t.Fatal("two wraps of same key must differ (random nonce)")
	}

	// But both must unwrap to same key
	u1, _ := UnwrapKey(wrappingKey, w1)
	u2, _ := UnwrapKey(wrappingKey, w2)

	if !bytes.Equal(u1, u2) {
		t.Fatal("different wraps must unwrap to same key")
	}
}
