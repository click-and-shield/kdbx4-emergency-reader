# List all entries and produce text output

## Command

```
keepass-reader.exe list .\doc\data\database.kdbx
```

> - You can use the sample database in the [`doc\data`](data) folder. The password is `dev`.
> - You may also specify the output format with the `--format` option: `keepass-reader.exe list .\doc\data\database.kdbx --format="text"`.
>   However, the default format is `text`.

## Result

```text
[1] /Password 1
[2] /password 2
[3] /group 1/Password 1
[4] /group 1/password 2
[5] /group 1/group 2/password 2
[6] /group 1/group 2/Password 1
[7] /group 3/Password 10
[8] /group 3/password 20
```

You then use the index (1, 2... 8) to select the entry.
