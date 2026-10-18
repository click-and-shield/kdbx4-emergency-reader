package reader

import "testing"

func TestEscapePathPart(t *testing.T) {
	if got, want := EscapePathPart(`a\b/c`), `a\\b\/c`; got != want {
		t.Fatalf("EscapePathPart() = %q, want %q", got, want)
	}
}
