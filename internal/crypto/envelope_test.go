package crypto

import (
	"bytes"
	"testing"
)

func TestEnvelope_WrapUnwrap(t *testing.T) {
	vaultKey, err := GenerateVaultKey(32)
	if err != nil {
		t.Fatal(err)
	}

	password := []byte("master-password")
	salt, _ := GenerateSalt(16)
	wrappingKey := DeriveKey(password, salt, DefaultKDFParams())

	wrapped, err := WrapKey(wrappingKey, vaultKey)
	if err != nil {
		t.Fatal(err)
	}

	unwrapped, err := UnwrapKey(wrappingKey, wrapped)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(vaultKey, unwrapped) {
		t.Fatal("unwrapped key doesn't match original")
	}
}

func TestEnvelope_WrongPassword(t *testing.T) {
	vaultKey, _ := GenerateVaultKey(32)
	salt, _ := GenerateSalt(16)
	p := DefaultKDFParams()

	wrappingKey1 := DeriveKey([]byte("correct"), salt, p)
	wrappingKey2 := DeriveKey([]byte("wrong"), salt, p)

	wrapped, _ := WrapKey(wrappingKey1, vaultKey)

	_, err := UnwrapKey(wrappingKey2, wrapped)
	if err == nil {
		t.Fatal("expected error with wrong password")
	}
}

func TestEnvelope_PasswordChange(t *testing.T) {
	vaultKey, _ := GenerateVaultKey(32)
	salt1, _ := GenerateSalt(16)
	p := DefaultKDFParams()

	oldKey := DeriveKey([]byte("old-password"), salt1, p)
	wrapped, _ := WrapKey(oldKey, vaultKey)

	unwrapped, err := UnwrapKey(oldKey, wrapped)
	if err != nil {
		t.Fatal(err)
	}

	salt2, _ := GenerateSalt(16)
	newKey := DeriveKey([]byte("new-password"), salt2, p)
	rewrapped, err := WrapKey(newKey, unwrapped)
	if err != nil {
		t.Fatal(err)
	}

	final, err := UnwrapKey(newKey, rewrapped)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(vaultKey, final) {
		t.Fatal("vault key should survive password change")
	}
}
