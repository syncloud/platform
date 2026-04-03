package main

import (
	"fmt"
	"os"
	"testapp/installer"
)

func main() {
	if err := installer.Configure(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
