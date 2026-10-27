package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func fixture() string { return filepath.Join("..", "..", "doc", "data", "database.kdbx") }

func TestHelp(t *testing.T) {
	var out bytes.Buffer
	if err := Run([]string{"help"}, "", &out); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"keepass-reader list <database>", "keepass-reader get <database> <id>",
		"keepass-reader help", "--format=text|json", "--user", "--password",
		"--url", "--notes", "--files", "--output=<dir>",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("help is missing %q", want)
		}
	}
}

func TestHelpRejectsArguments(t *testing.T) {
	var out bytes.Buffer
	if err := Run([]string{"help", "unexpected"}, "", &out); !errors.Is(err, ErrUsage) {
		t.Fatalf("expected usage error, got %v", err)
	}
	if out.Len() != 0 {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestHelpWriteError(t *testing.T) {
	want := errors.New("write failed")
	if err := Run([]string{"help"}, "", failingWriter{want}); !errors.Is(err, want) {
		t.Fatalf("expected write error, got %v", err)
	}
}

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

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
	var items []struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	}
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	wantLabels := []string{
		"/Password 1",
		"/password 2",
		"/group 1/Password 1",
		"/group 1/password 2",
		"/group 1/group 2/password 2",
		"/group 1/group 2/Password 1",
		"/group 3/Password 10",
		"/group 3/password 20",
	}
	if len(items) != len(wantLabels) {
		t.Fatalf("got %d items: %v", len(items), items)
	}
	for i, item := range items {
		if item.ID != i+1 || item.Label != wantLabels[i] {
			t.Errorf("item %d: got %+v, want id %d and label %q", i, item, i+1, wantLabels[i])
		}
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
