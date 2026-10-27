// Package keepass opens KeePass KDBX databases.
package keepass

import (
	"fmt"
	"io"
	"os"

	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

// Decrypt reads and decrypts a KDBX database protected only by password.
// Protected entry fields (notably passwords) are unlocked before the database
// is returned.
func Decrypt(r io.Reader, password string) (*gokeepasslib.Database, error) {
	if r == nil {
		return nil, fmt.Errorf("decrypt KDBX: nil reader")
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)
	if err := gokeepasslib.NewDecoder(r).Decode(db); err != nil {
		return nil, fmt.Errorf("decrypt KDBX: %w", err)
	}
	if err := db.UnlockProtectedEntries(); err != nil {
		return nil, fmt.Errorf("decrypt KDBX: %w", err)
	}
	return db, nil
}

// DecryptFile is the file-oriented counterpart of Decrypt.
func DecryptFile(path, password string) (*gokeepasslib.Database, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open KDBX %q: %w", path, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			_ = fmt.Errorf("close KDBX %q: %w", path, err)
		}
	}(f)

	return Decrypt(f, password)
}
