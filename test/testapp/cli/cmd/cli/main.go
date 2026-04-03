package main

import (
	"fmt"
	"os"
	"testapp/installer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: cli <command>")
		os.Exit(1)
	}
	var err error
	switch os.Args[1] {
	case "access-change":
		err = installer.AccessChange()
	case "storage-change":
		err = installer.StorageChange()
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
