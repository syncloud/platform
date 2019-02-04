package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func backups(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "[]")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.sock")
		return
	}

	os.Remove(os.Args[1])
	http.HandleFunc("/backup/list", backups)
	server := http.Server{}

	unixListener, err := net.Listen("unix", os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println("Started backend")
	server.Serve(unixListener)
}
