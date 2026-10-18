package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
	"keepass-reader/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: keepass-reader <list|get> ...")
		os.Exit(2)
	}

	fmt.Fprint(os.Stderr, "Mot de passe : ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "lecture du mot de passe:", err)
		os.Exit(1)
	}

	if err := cli.Run(os.Args[1:], string(password), os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "erreur:", err)
		if errors.Is(err, cli.ErrUsage) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
