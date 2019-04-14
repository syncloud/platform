package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/job"
)

type Backend struct {
	Master *job.Master
	backup *backup.Backup
	worker *job.Worker
}

func NewBackend(master *job.Master, backup *backup.Backup, worker *job.Worker) *Backend {
	return &Backend{
		Master: master,
		backup: backup,
		worker: worker,
	}
}

func (backend *Backend) Start(socket string) {
	go backend.worker.Start()
	http.HandleFunc("/job/status", Handle(backend.JobStatus))
	http.HandleFunc("/backup/list", Handle(backend.BackipList))
	http.HandleFunc("/backup/create", Handle(backend.BackupCreate))
	http.HandleFunc("/backup/restore", Handle(backend.BackupRestore))

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

func fail(w http.ResponseWriter, err error) {
	appError := err.Error()
	response := Response{
		Success: false,
		Message: &appError,
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
		fail(w, err)
	} else {
		fmt.Fprintf(w, string(responseJson))
	}
}

func Handle(f func(w http.ResponseWriter, req *http.Request) (interface{}, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("request: %s\n", req.URL.Path)
		data, err := f(w, req)
		if err != nil {
			fail(w, err)
		} else {
			success(w, data)
		}
	}
}

func (backend *Backend) BackipList(w http.ResponseWriter, req *http.Request) (interface{}, error) {
	return backend.backup.List()
}

func (backend *Backend) BackupCreate(w http.ResponseWriter, req *http.Request) (interface{}, error) {
	apps, ok := req.URL.Query()["app"]
	if !ok || len(apps) < 1 {
		return nil, errors.New("app is missing")
	}
	
	backend.Master.Offer(job.JobBackupCreate{App: apps[0], File: files[0]})
	return "submitted", nil
}

func (backend *Backend) BackupRestore(w http.ResponseWriter, req *http.Request) (interface{}, error) {
	apps, ok := req.URL.Query()["app"]
	if !ok || len(apps) < 1 {
		return nil, errors.New("app is missing")
	}
	files, ok := req.URL.Query()["file"]
	if !ok || len(files) < 1 {
		return nil, errors.New("file is missing")
	}

	backend.Master.Offer(job.JobBackupRestore{App: apps[0], File: files[0]})
	return "submitted", nil
}

func (backend *Backend) JobStatus(w http.ResponseWriter, req *http.Request) (interface{}, error) {
	return backend.Master.Status().String(), nil
}
