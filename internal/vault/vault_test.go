package vault

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	v := New()
	if v.Version != 1 {
		t.Fatalf("expected version 1, got %d", v.Version)
	}
	if len(v.Entries) != 0 {
		t.Fatal("new vault should have no entries")
	}
}

func TestAddRemoveEntry(t *testing.T) {
	v := New()
	e := Entry{
		ID:        "1",
		Title:     "Test",
		Username:  "user",
		Password:  "pass",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	v.AddEntry(e)

	if len(v.Entries) != 1 {
		t.Fatal("expected 1 entry")
	}

	if !v.RemoveEntry("1") {
		t.Fatal("should return true for existing entry")
	}
	if len(v.Entries) != 0 {
		t.Fatal("expected 0 entries after removal")
	}
	if v.RemoveEntry("1") {
		t.Fatal("should return false for missing entry")
	}
}

func TestFindEntry(t *testing.T) {
	v := New()
	v.AddEntry(Entry{ID: "1", Title: "Found"})

	if e := v.FindEntry("1"); e == nil || e.Title != "Found" {
		t.Fatal("expected to find entry")
	}
	if e := v.FindEntry("2"); e != nil {
		t.Fatal("expected nil for missing entry")
	}
}

func TestEntriesByCategory(t *testing.T) {
	v := New()
	v.AddEntry(Entry{ID: "1", Category: "work"})
	v.AddEntry(Entry{ID: "2", Category: "personal"})
	v.AddEntry(Entry{ID: "3", Category: "work"})

	work := v.EntriesByCategory("work")
	if len(work) != 2 {
		t.Fatalf("expected 2 work entries, got %d", len(work))
	}
}
