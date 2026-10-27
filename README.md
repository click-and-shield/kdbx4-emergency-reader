# KeepassXC database emergency reader

I wanted a portable standalone executable that I can use to read a KeePassXC database.

By _portable standalone executable_, I mean a zero-dependency executable that does not require any installer.

Such an executable does not exist. Thus, I decided to write one... with a little help of AI.

## Usage

```
$ keepass-reader.exe help                                       
Usage:
keepass-reader list <database> [--format=text|json]
keepass-reader get <database> <id> [options]
keepass-reader help

Commands:
list  List entries with their IDs and paths.
get   Show an entry by its ID from list.
help  Show this help without a database or password.

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

- [List all entries with their IDs and paths](doc/action-list-test.md).








```
keepass-reader.exe list .\doc\data\database.kdbx
```

Result:

```
[1] /Password 1
[2] /password 2
[3] /group 1/Password 1
[4] /group 1/password 2
[5] /group 1/group 2/password 2
[6] /group 1/group 2/Password 1
[7] /group 3/Password 10
[8] /group 3/password 20
```

```
keepass-reader.exe list .\doc\data\database.kdbx --format=json
```

Result:

```
[
  "/Password 1",
  "/password 2",
  "/group 1/Password 1",
  "/group 1/password 2",
  "/group 1/group 2/password 2",
  "/group 1/group 2/Password 1",
  "/group 3/Password 10",
  "/group 3/password 20"
]
```

```
keepass-reader.exe get .\doc\data\database.kdbx 1 --user --password
```

Result:

```
user: user
password: passwod
```

```
keepass-reader.exe get .\doc\data\database.kdbx 1 --user --password --notes --url
```

Result:

```
user: user
password: passwod
url: https://example.com
notes:
This is a note.
```

```
keepass-reader.exe get .\doc\data\database.kdbx 1 --user --password --notes --url --format=json
```

Result:

```
{
  "notes": "This is a note.",
  "password": "passwod",
  "url": "https://example.com",
  "user": "user"
}
```

```
keepass-reader.exe get .\doc\data\database.kdbx 2 --output=.\data 
```

Result:

```
user: user
password: password
url: https://exemple2.com
notes:
This is password 2.
files: "doc-text.txt", "meetup.png"
```

```
C:\keepass-reader>dir data
 Le volume dans le lecteur C n’a pas de nom.
 Le numéro de série du volume est 608C-9134

 Répertoire de C:\Users\denis\Documents\github\keepass-reader\data

18/10/2026  21:41    <DIR>          .
18/10/2026  21:42    <DIR>          ..
18/10/2026  21:41                14 doc-text.txt
18/10/2026  21:41           141 145 meetup.png
               2 fichier(s)          141 159 octets
               2 Rép(s)  438 130 188 88 octets libres
```
