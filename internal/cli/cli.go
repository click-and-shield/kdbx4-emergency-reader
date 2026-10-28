// Package cli implements the keepass-reader command-line interface.
// Parse the command line and execute the requested action ("list", "get", or "help").
package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"

	"keepass-reader/internal/keepass"
	"keepass-reader/internal/reader"

	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

var ErrUsage = errors.New("invalid command line")

var version = "undeclared version"

var buildDate = "unknown compilation date"

type getOptions struct {
	user, password, url, notes, files bool
	output, format                    string
}

// Run parses the command line and executes the requested action ("list", "get", or "help").
// Note: Password acquisition is kept outside this function so the command logic can be tested without a terminal.
func Run(args []string, password string, stdout io.Writer) error {
	if len(args) == 0 {
		return usageError("expected command list, get, or help")
	}
	switch args[0] {
	case "help":
		if len(args) != 1 {
			return usageError("usage: keepass-reader help")
		}
		return printHelp(stdout)
	case "version":
		_, err := io.WriteString(stdout, fmt.Sprintf("%s (%s-%s %s)\n", version, runtime.GOOS, runtime.GOARCH, buildDate))
		return err
	case "list":
		return runList(args[1:], password, stdout)
	case "get":
		return runGet(args[1:], password, stdout)
	default:
		return usageError("unknown command %q", args[0])
	}
}

func printHelp(out io.Writer) error {
	_, err := io.WriteString(out, fmt.Sprintf(`
version: %s (%s-%s %s)

Usage:
  keepass-reader list <database> [--format=text|json]
  keepass-reader get <database> <id> [options]
  keepass-reader help
  keepass-reader version

Commands:
  list    List entries with their IDs and paths.
  get     Show an entry by its ID from list.
  help    Show this help without a database or password.
  version Show the version.

Options for list and get:
  --format=text|json  Output format (default: text).

Options for get:
  --user             Show username.
  --password         Show password.
  --url              Show URL.
  --notes            Show notes.
  --files            Save attachments.
  --output=<dir>     Attachment output directory (default: current directory).

If no fields are selected for get, all fields are shown and attachments are saved.
`, version, runtime.GOOS, runtime.GOARCH, buildDate))

	return err
}

func runList(args []string, password string, out io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	format := fs.String("format", "text", "text or json")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"--format": true})); err != nil {
		return usageError("list: %v", err)
	}
	if fs.NArg() != 1 {
		return usageError("usage: keepass-reader list <database> [--format=text|json]")
	}
	if err := validFormat(*format); err != nil {
		return err
	}
	db, err := keepass.DecryptFile(fs.Arg(0), password)
	if err != nil {
		return err
	}
	entries := reader.Entries(db)
	if *format == "json" {
		items := make([]struct {
			ID    int    `json:"id"`
			Label string `json:"label"`
		}, len(entries))
		for i := range entries {
			items[i].ID = i + 1
			items[i].Label = entries[i].Path
		}
		return writeJSON(out, items)
	}
	for i := range entries {
		if _, err := fmt.Fprintf(out, "[%d] %s\n", i+1, entries[i].Path); err != nil {
			return err
		}
	}
	return nil
}

func runGet(args []string, password string, out io.Writer) error {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	opts := getOptions{}
	fs.BoolVar(&opts.user, "user", false, "show username")
	fs.BoolVar(&opts.password, "password", false, "show password")
	fs.BoolVar(&opts.url, "url", false, "show URL")
	fs.BoolVar(&opts.notes, "notes", false, "show notes")
	fs.BoolVar(&opts.files, "files", false, "save attachments")
	fs.StringVar(&opts.output, "output", ".", "attachment output directory")
	fs.StringVar(&opts.format, "format", "text", "text or json")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"--format": true, "--output": true})); err != nil {
		return usageError("get: %v", err)
	}
	if fs.NArg() != 2 {
		return usageError("usage: keepass-reader get <database> <id> [options]")
	}
	if err := validFormat(opts.format); err != nil {
		return err
	}
	id, err := strconv.Atoi(fs.Arg(1))
	if err != nil || id < 1 {
		return usageError("id must be a positive integer")
	}
	selected := opts.user || opts.password || opts.url || opts.notes || opts.files
	if !selected {
		opts.user, opts.password, opts.url, opts.notes, opts.files = true, true, true, true, true
	}

	db, err := keepass.DecryptFile(fs.Arg(0), password)
	if err != nil {
		return err
	}
	entries := reader.Entries(db)
	if id > len(entries) {
		return fmt.Errorf("entry id %d does not exist (valid range: 1-%d)", id, len(entries))
	}
	entry := entries[id-1].Value
	if opts.files {
		if err := reader.SaveAttachments(db, entry, opts.output); err != nil {
			return err
		}
	}
	if opts.format == "json" {
		return renderJSON(out, entry, opts)
	}
	return renderText(out, entry, opts)
}

func renderText(out io.Writer, entry *gokeepasslib.Entry, opts getOptions) error {
	var b strings.Builder
	var err error
	if opts.user {
		if _, err = fmt.Fprintf(&b, "user: %s\n", entry.GetContent("UserName")); err != nil {
			return err
		}
	}
	if opts.password {
		if _, err = fmt.Fprintf(&b, "password: %s\n", entry.GetContent("Password")); err != nil {
			return err
		}
	}
	if opts.url {
		if _, err = fmt.Fprintf(&b, "url: %s\n", entry.GetContent("URL")); err != nil {
			return err
		}
	}
	if opts.notes {
		if _, err = fmt.Fprintf(&b, "notes:\n%s\n", entry.GetContent("Notes")); err != nil {
			return err
		}
	}
	if opts.files {
		if _, err = fmt.Fprintf(&b, "files: %s\n", quoteNames(reader.AttachmentNames(entry))); err != nil {
			return err
		}
	}
	_, err = io.WriteString(out, b.String())
	return err
}

func renderJSON(out io.Writer, entry *gokeepasslib.Entry, opts getOptions) error {
	data := make(map[string]any)
	set := func(key, value string) {
		if value == "" {
			data[key] = nil
		} else {
			data[key] = value
		}
	}
	if opts.user {
		set("user", entry.GetContent("UserName"))
	}
	if opts.password {
		set("password", entry.GetContent("Password"))
	}
	if opts.url {
		set("url", entry.GetContent("URL"))
	}
	if opts.notes {
		set("notes", entry.GetContent("Notes"))
	}
	if opts.files {
		data["files"] = reader.AttachmentNames(entry)
	}
	return writeJSON(out, data)
}

func validFormat(format string) error {
	if format != "text" && format != "json" {
		return usageError("unsupported format %q (expected text or json)", format)
	}
	return nil
}

func writeJSON(out io.Writer, value any) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

func quoteNames(names []string) string {
	quoted := make([]string, len(names))
	for i := range names {
		quoted[i] = strconv.Quote(names[i])
	}
	return strings.Join(quoted, ", ")
}

// flagsFirst allows the documented "positional arguments, then flags" form,
// which the standard flag package does not accept by itself.
func flagsFirst(args []string, valueFlags map[string]bool) []string {
	var flags, positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") {
			positional = append(positional, arg)
			continue
		}
		flags = append(flags, arg)
		if valueFlags[arg] && i+1 < len(args) {
			i++
			flags = append(flags, args[i])
		}
	}
	return append(flags, positional...)
}

func usageError(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrUsage, fmt.Sprintf(format, args...))
}
