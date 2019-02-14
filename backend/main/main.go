package main

import (
	"log"
	"os"

	"github.com/syncloud/platform/rest"

)

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: ", os.Args[0], "/path.sock")
		return
	}

	os.Remove(os.Args[1])
 backend := rest.NewBackend()
  backend.Start(os.Args[1])

}
