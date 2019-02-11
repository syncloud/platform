package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/syncloud/platform/backup"
)

type Response struct {
	Success bool         `json:"success"`
	Message *string      `json:"message,omitempty"`
	Data    *interface{} `json:"data,omitempty"`
}

func fail(w http.ResponseWriter, err error, message string) {
	log.Println(err)
	response := Response{
		Success: false,
		Message: &message,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, string(responseJson))
	}
}

func success(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    &data,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		fail(w, err, "Cannot encode to JSON")
	} else {
		fmt.Println(response.Success)
		fmt.Fprintf(w, string(responseJson))
	}
}

func Handle(f func() (interface{}, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := f()
		if err != nil {
			fail(w, err, "Cannot get data")
		} else {
			success(w, data)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: ", os.Args[0], "/path.sock")
		return
	}

	os.Remove(os.Args[1])
	backup := backup.NewDefault()
	http.HandleFunc("/backup/list", Handle(func() (interface{}, error) { return backup.List() }))
	server := http.Server{}

	unixListener, err := net.Listen("unix", os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Println("Started backend")
	server.Serve(unixListener)
}
