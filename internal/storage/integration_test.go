package storage

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/walteraandrade/cerberus/internal/crypto"
	"github.com/walteraandrade/cerberus/internal/vault"
)

func TestFullVaultLifecycle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	password := []byte("integration-test-password")
	params := crypto.DefaultKDFParams()

	// 1. Create vault with initial entry
	v := vault.New()
	v.AddEntry(vault.Entry{
		ID:        "entry-1",
		Title:     "GitHub",
		URL:       "https://github.com",
		Username:  "testuser",
		Password:  "gh-secret-123",
		Category:  "dev",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	data, err := vault.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := CreateVault(path, password, data, params); err != nil {
		t.Fatalf("create: %v", err)
	}

	// 2. Reopen and verify
	pt, err := OpenVault(path, password, params)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	v2, err := vault.Unmarshal(pt)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(v2.Entries) != 1 || v2.Entries[0].Title != "GitHub" {
		t.Fatal("initial entry missing after reopen")
	}

	// 3. Add second entry and save
	v2.AddEntry(vault.Entry{
		ID:        "entry-2",
		Title:     "Gmail",
		Username:  "test@gmail.com",
		Password:  "gmail-pass-456",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	data, _ = vault.Marshal(v2)
	if err := CreateVault(path, password, data, params); err != nil {
		t.Fatalf("save after add: %v", err)
	}

	// 4. Edit first entry
	pt, _ = OpenVault(path, password, params)
	v3, _ := vault.Unmarshal(pt)

	entry := v3.FindEntry("entry-1")
	if entry == nil {
		t.Fatal("entry-1 not found")
	}
	entry.Password = "updated-password"
	entry.UpdatedAt = time.Now()

	data, _ = vault.Marshal(v3)
	if err := CreateVault(path, password, data, params); err != nil {
		t.Fatalf("save after edit: %v", err)
	}

	// 5. Verify edit persisted
	pt, _ = OpenVault(path, password, params)
	v4, _ := vault.Unmarshal(pt)
	edited := v4.FindEntry("entry-1")
	if edited.Password != "updated-password" {
		t.Fatal("edit not persisted")
	}
	if len(v4.Entries) != 2 {
		t.Fatal("expected 2 entries after edit")
	}

	// 6. Delete entry
	v4.RemoveEntry("entry-2")
	data, _ = vault.Marshal(v4)
	if err := CreateVault(path, password, data, params); err != nil {
		t.Fatalf("save after delete: %v", err)
	}

	// 7. Final reopen — verify delete
	pt, _ = OpenVault(path, password, params)
	v5, _ := vault.Unmarshal(pt)
	if len(v5.Entries) != 1 {
		t.Fatalf("expected 1 entry after delete, got %d", len(v5.Entries))
	}
	if v5.Entries[0].ID != "entry-1" {
		t.Fatal("wrong entry survived delete")
	}
}

func TestPasswordChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	oldPw := []byte("old-password")
	newPw := []byte("new-password")
	params := crypto.DefaultKDFParams()

	v := vault.New()
	v.AddEntry(vault.Entry{ID: "1", Title: "Test", Password: "secret"})
	data, _ := vault.Marshal(v)

	if err := CreateVault(path, oldPw, data, params); err != nil {
		t.Fatal(err)
	}

	// Open with old password
	pt, err := OpenVault(path, oldPw, params)
	if err != nil {
		t.Fatal(err)
	}

	// Re-create with new password (simulates password change)
	if err := CreateVault(path, newPw, pt, params); err != nil {
		t.Fatal(err)
	}

	// Old password should fail
	_, err = OpenVault(path, oldPw, params)
	if err == nil {
		t.Fatal("old password should fail")
	}

	// New password should work
	pt2, err := OpenVault(path, newPw, params)
	if err != nil {
		t.Fatal(err)
	}
	v2, _ := vault.Unmarshal(pt2)
	if len(v2.Entries) != 1 || v2.Entries[0].Password != "secret" {
		t.Fatal("data should survive password change")
	}
}

func TestCorruptedVault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	password := []byte("password")
	params := crypto.DefaultKDFParams()

	v := vault.New()
	data, _ := vault.Marshal(v)
	if err := CreateVault(path, password, data, params); err != nil {
		t.Fatal(err)
	}

	// Corrupt the file
	raw, _ := Read(path)
	raw.Ciphertext[0] ^= 0xFF
	if err := Write(path, raw.Salt, raw.WrappedKey, raw.Ciphertext); err != nil {
		t.Fatal(err)
	}

	_, err := OpenVault(path, password, params)
	if err == nil {
		t.Fatal("corrupted vault should fail to open")
	}
}

func TestEmptyVaultRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	password := []byte("password")
	params := crypto.DefaultKDFParams()

	v := vault.New()
	data, _ := vault.Marshal(v)
	if err := CreateVault(path, password, data, params); err != nil {
		t.Fatal(err)
	}

	pt, err := OpenVault(path, password, params)
	if err != nil {
		t.Fatal(err)
	}
	v2, _ := vault.Unmarshal(pt)
	if v2.Version != 1 || len(v2.Entries) != 0 {
		t.Fatal("empty vault round-trip failed")
	}
}
