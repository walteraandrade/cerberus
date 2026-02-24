package vault

import (
	"testing"
	"time"
)

func TestMarshalUnmarshal(t *testing.T) {
	v := New()
	v.AddEntry(Entry{
		ID:        "1",
		Title:     "GitHub",
		URL:       "https://github.com",
		Username:  "user",
		Password:  "secret",
		Category:  "dev",
		CreatedAt: time.Now().Truncate(time.Second),
		UpdatedAt: time.Now().Truncate(time.Second),
	})
	v.AddCategory(Category{ID: "dev", Name: "Development"})

	data, err := Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	v2, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}

	if v2.Version != v.Version {
		t.Fatal("version mismatch")
	}
	if len(v2.Entries) != 1 {
		t.Fatal("expected 1 entry")
	}
	if v2.Entries[0].Title != "GitHub" {
		t.Fatal("title mismatch")
	}
	if len(v2.Categories) != 1 {
		t.Fatal("expected 1 category")
	}
}
