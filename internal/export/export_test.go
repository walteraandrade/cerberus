package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/walteraandrade/cerberus/internal/vault"
)

func TestJSONExporter(t *testing.T) {
	v := vault.New()
	v.AddEntry(vault.Entry{ID: "1", Title: "GitHub", Username: "user", Password: "secret", Category: "dev"})
	v.AddEntry(vault.Entry{ID: "2", Title: "Gmail", Username: "user@gmail.com", Password: "pass"})

	data, err := JSONExporter{}.Export(v)
	if err != nil {
		t.Fatal(err)
	}

	var result exportData
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	if result.Count != 2 {
		t.Fatalf("expected 2 entries, got %d", result.Count)
	}
	if result.Entries[0].Title != "GitHub" {
		t.Fatal("expected GitHub entry first")
	}
	if result.ExportedAt == "" {
		t.Fatal("expected exported_at timestamp")
	}
}

func TestToFile(t *testing.T) {
	v := vault.New()
	v.AddEntry(vault.Entry{ID: "1", Title: "Test", Password: "pw"})

	dir := t.TempDir()
	path := filepath.Join(dir, "export.json")

	if err := ToFile(path, JSONExporter{}, v); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var result exportData
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result.Count != 1 {
		t.Fatalf("expected 1 entry, got %d", result.Count)
	}
}
