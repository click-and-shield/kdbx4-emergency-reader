package cli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func fixture() string { return filepath.Join("..", "doc", "data", "database.kdbx") }

func TestListText(t *testing.T) {
	var out bytes.Buffer
	if err := Run([]string{"list", fixture()}, "dev", &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 8 {
		t.Fatalf("got %d entries:\n%s", len(lines), out.String())
	}
	if !strings.HasPrefix(lines[0], "[1] /") || !strings.HasPrefix(lines[7], "[8] /") {
		t.Fatalf("unexpected numbering:\n%s", out.String())
	}
}

func TestListJSONWithTrailingFlag(t *testing.T) {
	var out bytes.Buffer
	if err := Run([]string{"list", fixture(), "--format=json"}, "dev", &out); err != nil {
		t.Fatal(err)
	}
	var paths []string
	if err := json.Unmarshal(out.Bytes(), &paths); err != nil {
		t.Fatal(err)
	}
	if len(paths) != 8 {
		t.Fatalf("got %d paths: %v", len(paths), paths)
	}
}

func TestGetSelectedJSONFields(t *testing.T) {
	var out bytes.Buffer
	if err := Run([]string{"get", fixture(), "1", "--user", "--format", "json"}, "dev", &out); err != nil {
		t.Fatal(err)
	}
	var data map[string]any
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	if len(data) != 1 {
		t.Fatalf("expected one selected field, got %v", data)
	}
	if _, ok := data["user"]; !ok {
		t.Fatalf("missing user field: %v", data)
	}
}

func TestGetRejectsUnknownID(t *testing.T) {
	err := Run([]string{"get", fixture(), "999", "--user"}, "dev", &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("unexpected error: %v", err)
	}
}
