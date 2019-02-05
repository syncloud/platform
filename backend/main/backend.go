package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"log"
	"encoding/json"

	"github.com/syncloud/platform/backup"
)

const backupDir = "/data/platform/backup"

func backups(w http.ResponseWriter, req *http.Request) {
	files, err := backup.List(backupDir)
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "")
		return
	}
	
	filesJson, err := json.Marshal(files)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
		fmt.Fprintf(w, "")
		return
	}
	fmt.Fprintf(w, string(filesJson))
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
