package export

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/walteraandrade/cerberus/internal/vault"
)

type Exporter interface {
	Export(v *vault.Vault) ([]byte, error)
}

type JSONExporter struct{}

type exportEntry struct {
	Title    string `json:"title"`
	URL      string `json:"url,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Notes    string `json:"notes,omitempty"`
	Category string `json:"category,omitempty"`
}

type exportData struct {
	ExportedAt string        `json:"exported_at"`
	Count      int           `json:"count"`
	Entries    []exportEntry `json:"entries"`
}

func (JSONExporter) Export(v *vault.Vault) ([]byte, error) {
	entries := make([]exportEntry, len(v.Entries))
	for i, e := range v.Entries {
		entries[i] = exportEntry{
			Title:    e.Title,
			URL:      e.URL,
			Username: e.Username,
			Password: e.Password,
			Notes:    e.Notes,
			Category: e.Category,
		}
	}

	data := exportData{
		ExportedAt: time.Now().Format(time.RFC3339),
		Count:      len(entries),
		Entries:    entries,
	}

	return json.MarshalIndent(data, "", "  ")
}

func ToFile(path string, exporter Exporter, v *vault.Vault) error {
	data, err := exporter.Export(v)
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
