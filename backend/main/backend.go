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

type Response struct {
	success bool `json:",omitempty"`
	message string `json:",omitempty"`
	data interface{} `json:",omitempty"`
}

func fail(w http.ResponseWriter, err error, message string) {
	log.Println(err)
	response := Response {
		success: false,
		message: message,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, string(responseJson))
	}
}

func success(w http.ResponseWriter, data interface{}) {
	response := Response {
		success: true,
		data: data,
	}
	reaponseJson, err := json.Marshal(response)
	if err != nil {
		fail(w, err, "Cannot encode to JSON")
	} else {
		fmt.Fprintf(w, string(reaponseJson))
	}
}

func backups(w http.ResponseWriter, req *http.Request) {
	files, err := backup.ListDefault()
	if err != nil {
		fail(w, err, "Cannot get list of backups")
	} else {
		success(w, files)
	}
	
}

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: ", os.Args[0], "/path.sock")
		return
	}

	os.Remove(os.Args[1])
	http.HandleFunc("/backup/list", backups)
	server := http.Server{}

	unixListener, err := net.Listen("unix", os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Println("Started backend")
	server.Serve(unixListener)
}
