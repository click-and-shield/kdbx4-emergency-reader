package main

import (
	"fmt"
	"log"
	"os"

	"keepass-reader/keepass"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s DATABASE.kdbx PASSWORD", os.Args[0])
	}

	db, err := keepass.DecryptFile(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Base déchiffrée : %s\n", db.Content.Meta.DatabaseName)
}
