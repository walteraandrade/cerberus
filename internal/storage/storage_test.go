package storage

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/walteraandrade/cerberus/internal/crypto"
	"github.com/walteraandrade/cerberus/internal/vault"
)

func TestWriteRead_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.enc")

	salt := []byte("0123456789abcdef")
	wrappedKey := []byte("wrapped-key-data-here")
	ciphertext := []byte("encrypted-vault-data")

	if err := Write(path, salt, wrappedKey, ciphertext); err != nil {
		t.Fatal(err)
	}

	vf, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(vf.Salt, salt) {
		t.Fatal("salt mismatch")
	}
	if !bytes.Equal(vf.WrappedKey, wrappedKey) {
		t.Fatal("wrapped key mismatch")
	}
	if !bytes.Equal(vf.Ciphertext, ciphertext) {
		t.Fatal("ciphertext mismatch")
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.enc")

	if Exists(path) {
		t.Fatal("should not exist yet")
	}

	os.WriteFile(path, []byte("data"), 0600)
	if !Exists(path) {
		t.Fatal("should exist after write")
	}
}

func TestCreateOpenVault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")

	v := vault.New()
	v.AddEntry(vault.Entry{ID: "1", Title: "Test", Password: "secret"})

	data, err := vault.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	password := []byte("master-password")
	params := crypto.DefaultKDFParams()

	if err := CreateVault(path, password, data, params); err != nil {
		t.Fatal(err)
	}

	plaintext, err := OpenVault(path, password, params)
	if err != nil {
		t.Fatal(err)
	}

	v2, err := vault.Unmarshal(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	if len(v2.Entries) != 1 || v2.Entries[0].Title != "Test" {
		t.Fatal("vault data mismatch after open")
	}
}

func TestOpenVault_WrongPassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")

	params := crypto.DefaultKDFParams()
	if err := CreateVault(path, []byte("correct"), []byte("data"), params); err != nil {
		t.Fatal(err)
	}

	_, err := OpenVault(path, []byte("wrong"), params)
	if err == nil {
		t.Fatal("expected error with wrong password")
	}
}

func TestRead_InvalidMagic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.enc")
	os.WriteFile(path, []byte("BADMAGIC"), 0600)

	_, err := Read(path)
	if err == nil {
		t.Fatal("expected error for invalid magic")
	}
}
