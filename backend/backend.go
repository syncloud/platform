package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "request %s", req.URL.Path[1:])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.sock")
		return
	}

	fmt.Println("Starting backend")

	os.Remove(os.Args[1])
	http.HandleFunc("/", hello)
	server := http.Server{}

	unixListener, err := net.Listen("unix", os.Args[1])
	if err != nil {
		panic(err)
	}
	server.Serve(unixListener)
}
