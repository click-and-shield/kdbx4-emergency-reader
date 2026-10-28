# KeepassXC database emergency reader

I wanted a portable standalone executable that I can use to read a KeePassXC database.

By _portable standalone executable_, I mean a zero-dependency executable that does not require any installer.

Such an executable does not exist. Thus, I decided to write one... with a little help of AI.

> **WARNING**: this application is intended to be used in emergency situations.
> Typically:
> - You only have your database file and this application on a USB key.
> - You cannot use your own computer, nor your smartphone, to read the database.
>
> You have no choice but to use a _potentially_ insecure device.
> The risk is that someone steals your database file and your password.

## Usage

```
$ keepass-reader.exe help                                       
version: 1.0.0 (windows-amd64)

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
```

Examples:

- [List all entries with their IDs and paths](doc/action-list.md)
- [Get the data associated with a given entry identified by its ID](doc/action-get.md)

## Compilation

- **PowerShell**:

  ```powershell
  $buildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-dd'T'HH:mm:ss'Z'")
  $version = "1.0.0"
  echo $buildDate
  echo $version
  go build -ldflags "-X keepass-reader/internal/cli.buildDate=$buildDate -X keepass-reader/internal/cli.version=$version" -o keepass-reader.exe .
  ```

- **Bash**:

  ```bash
  go build -ldflags "-X keepass-reader/internal/cli.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X keepass-reader/internal/cli.version=$version" -o keepass-reader.exe .
  ```

## Execute unit tests

  ```
  go test -v -cover ./...
  ```
