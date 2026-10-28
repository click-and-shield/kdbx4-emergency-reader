# Get entry's data

## Output text

### Command

#### Get all data

```
keepass-reader.exe get .\doc\data\database.kdbx 1
```

> This is equivalent to:
> 
> ```
> keepass-reader.exe get .\doc\data\database.kdbx 1  --user --password --url --notes --files
> ```

#### Result

```
user: user
password: passwod
url: https://example.com
notes:
This is a note.
files: 
```



#### Get only the user and password

```
keepass-reader.exe get .\doc\data\database.kdbx 1 --user --password
```

#### Result

```
user: user
password: passwod
```


