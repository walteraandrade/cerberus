package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "test-key-32-bytes-long-padding!!")

	plaintext := []byte("secret data")
	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatal("decrypted data doesn't match original")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	copy(key1, "key-one-32-bytes-long-padding!!!")
	copy(key2, "key-two-32-bytes-long-padding!!!")

	ciphertext, err := Encrypt(key1, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = Decrypt(key2, ciphertext)
	if err == nil {
		t.Fatal("expected error decrypting with wrong key")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "test-key-32-bytes-long-padding!!")

	ciphertext, err := Encrypt(key, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err = Decrypt(key, ciphertext)
	if err == nil {
		t.Fatal("expected error for tampered ciphertext")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	key := make([]byte, 32)
	_, err := Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}
