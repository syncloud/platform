package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/syncloud/platform/backup"
 "github.com/syncloud/platform/job"
)

type Backend struct {
 Master *job.Master
 Backup *backup.Backup
 worker *job.Worker
}

func NewBackend() *Backend {
 master := job.NewMaster()
	backup := backup.NewDefault()
 worker := job.NewWorker(master.JobQueue(), backup)

 return &Backend{
  Master: master,
	Backup: backup,
 worker: worker,
 }
}

func (backend *Backend) Start(socket string) {
 backend.worker.Start()
	http.HandleFunc("/backup/list", Handle(backend.BackipList))
 http.HandleFunc("/backup/create", Handle(backend.BackupCreate))
 
	server := http.Server{}

	unixListener, err := net.Listen("unix", socket)
	if err != nil {
		panic(err)
	}
	log.Println("Started backend")
	server.Serve(unixListener)

}

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

func Handle(f func(w http.ResponseWriter, req *http.Request) (interface{}, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := f(w, req)
		if err != nil {
			fail(w, err, "Cannot get data")
		} else {
			success(w, data)
		}
	}
}

func (backend *Backend) BackipList(w http.ResponseWriter, req *http.Request) (interface{}, error) {
	v,e := backend.Backup.List()
 return v, e
}

func (backend *Backend) BackupCreate(w http.ResponseWriter, req *http.Request) (interface{}, error) {
 apps, ok := req.URL.Query()["app"]
 if !ok || len(apps) < 1 {
  return "app is missing", nil
 }
 files, ok := req.URL.Query()["file"]
 if !ok || len(files) <1 {
  return "file is missing", nil
 }

	backend.Master.BackupCreateJob(apps[0], files[0])
	return "submitted", nil
}
