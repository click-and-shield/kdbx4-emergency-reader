package main

import (
	"errors"
	"fmt"
	"os"

	"keepass-reader/internal/cli"

	"golang.org/x/term"
)

func main() {
	// Check the command line
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "usage: keepass-reader <list|get> ...")
		os.Exit(2)
	}

	// Read the password
	_, _ = fmt.Fprint(os.Stderr, "Password : ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	_, _ = fmt.Fprintln(os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "An error occurred while reading the password:", err)
		os.Exit(1)
	}

	// Run the command
	if err := cli.Run(os.Args[1:], string(password), os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "erreur:", err)
		if errors.Is(err, cli.ErrUsage) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
