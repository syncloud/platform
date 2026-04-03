package main

import (
	"fmt"
	"os"
	"testapp/installer"
)

func main() {
	if err := installer.Install(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
