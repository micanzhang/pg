package pg

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFileEntryPersistant(t *testing.T) {
	path := fmt.Sprintf("pg-%d-fep", time.Now().Unix())
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			t.Error(err)
		}
	}

	defer os.Remove(path)

	e := NewFileEntryPersistant(path)

	entries := []Entry{
		Entry{
			Domain:   "google.com",
			Username: "allen",
			Password: Password("allen@Admin"),
		},
		Entry{
			Domain:   "mysql",
			Username: "root",
			Password: Password(""),
		},
		Entry{
			Domain:   "",
			Username: "",
			Password: Password("123!@#QWRsasd"),
		},
	}

	err := e.Save(entries)
	if err != nil {
		t.Fatal(err)
	}

	entries1, err := e.Restore()
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != len(entries1) {
		t.Fatalf("expected: %d, got: %d", len(entries), len(entries1))
	}

	for i, entry := range entries {
		entry1 := entries1[i]

		if entry.Domain != entry1.Domain {
			t.Errorf("expected: %s, got: %s", entry.Domain, entry1.Domain)
		}

		if entry.Username != entry1.Username {
			t.Errorf("expected: %s, got: %s", entry.Username, entry1.Username)
		}

		if entry.Password != entry1.Password {
			t.Errorf("expected: %s, got: %s", entry.Password, entry1.Password)
		}
	}
}
