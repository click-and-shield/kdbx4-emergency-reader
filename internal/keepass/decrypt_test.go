package keepass

import (
	"path/filepath"
	"testing"

	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

func TestDecryptFile(t *testing.T) {
	path := filepath.Join("..", "..", "doc", "data", "database.kdbx")
	db, err := DecryptFile(path, "dev")
	if err != nil {
		t.Fatalf("DecryptFile: %v", err)
	}

	titles := collectTitles(db.Content.Root.Groups)
	if !titles["Password 1"] || !titles["password 2"] {
		t.Fatalf("expected exported fixture entries, got titles %v", titles)
	}
}

func TestDecryptFileRejectsWrongPassword(t *testing.T) {
	path := filepath.Join("..", "..", "doc", "data", "database.kdbx")
	if _, err := DecryptFile(path, "wrong"); err == nil {
		t.Fatal("DecryptFile accepted an invalid password")
	}
}

func collectTitles(groups []gokeepasslib.Group) map[string]bool {
	titles := make(map[string]bool)
	var visit func([]gokeepasslib.Group)
	visit = func(groups []gokeepasslib.Group) {
		for _, group := range groups {
			for _, entry := range group.Entries {
				titles[entry.GetTitle()] = true
			}
			visit(group.Groups)
		}
	}
	visit(groups)
	return titles
}
