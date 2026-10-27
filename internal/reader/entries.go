// Package reader exposes the data needed by the command-line interface.
package reader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

// Entry is an entry together with its display path.
type Entry struct {
	Path  string
	Value *gokeepasslib.Entry
}

// Entries returns entries in KeePass display order. The database root group's
// name is deliberately omitted from paths.
func Entries(db *gokeepasslib.Database) []Entry {
	var result []Entry
	for i := range db.Content.Root.Groups {
		root := &db.Content.Root.Groups[i]
		appendGroupEntries(&result, root, "", false)
	}
	return result
}

func appendGroupEntries(dst *[]Entry, group *gokeepasslib.Group, parent string, includeName bool) {
	path := parent
	if includeName {
		path += "/" + EscapePathPart(group.Name)
	}
	for i := range group.Entries {
		entry := &group.Entries[i]
		*dst = append(*dst, Entry{Path: path + "/" + EscapePathPart(entry.GetTitle()), Value: entry})
	}
	for i := range group.Groups {
		appendGroupEntries(dst, &group.Groups[i], path, true)
	}
}

// EscapePathPart escapes the two separators used in displayed entry paths.
func EscapePathPart(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	return strings.ReplaceAll(value, "/", `\/`)
}

// AttachmentNames returns attachment names in database order.
func AttachmentNames(entry *gokeepasslib.Entry) []string {
	names := make([]string, len(entry.Binaries))
	for i := range entry.Binaries {
		names[i] = entry.Binaries[i].Name
	}
	return names
}

// SaveAttachments writes all attachments without allowing path traversal or
// overwriting an existing file.
func SaveAttachments(db *gokeepasslib.Database, entry *gokeepasslib.Entry, output string) error {
	if err := os.MkdirAll(output, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	for i := range entry.Binaries {
		ref := &entry.Binaries[i]
		if ref.Name == "" || filepath.IsAbs(ref.Name) || filepath.Base(ref.Name) != ref.Name {
			return fmt.Errorf("unsafe attachment name %q", ref.Name)
		}
		binary := ref.Find(db)
		if binary == nil {
			return fmt.Errorf("attachment %q has no binary content", ref.Name)
		}
		data, err := binary.GetContentBytes()
		if err != nil {
			return fmt.Errorf("decode attachment %q: %w", ref.Name, err)
		}
		path := filepath.Join(output, ref.Name)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return fmt.Errorf("create attachment %q: %w", path, err)
		}
		_, writeErr := file.Write(data)
		closeErr := file.Close()
		if writeErr != nil {
			return fmt.Errorf("write attachment %q: %w", path, writeErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close attachment %q: %w", path, closeErr)
		}
	}
	return nil
}
