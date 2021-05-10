package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest/model"
	"log"
	"net"
	"net/http"

	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/job"
)

type Backend struct {
	Master       *job.Master
	backup       *backup.Backup
	eventTrigger *event.Trigger
	worker       *job.Worker
	redirect     *redirect.Redirect
}

func NewBackend(master *job.Master, backup *backup.Backup, eventTrigger *event.Trigger, worker *job.Worker, redirect *redirect.Redirect) *Backend {
	return &Backend{
		Master:       master,
		backup:       backup,
		eventTrigger: eventTrigger,
		worker:       worker,
		redirect:     redirect,
	}
}

func (backend *Backend) Start(network string, address string) {
	unixListener, err := net.Listen(network, address)
	if err != nil {
		panic(err)
	}

	go backend.worker.Start()
	http.HandleFunc("/job/status", Handle(http.MethodGet, backend.JobStatus))
	http.HandleFunc("/backup/list", Handle(http.MethodGet, backend.BackupList))
	http.HandleFunc("/backup/create", Handle(http.MethodPost, backend.BackupCreate))
	http.HandleFunc("/backup/restore", Handle(http.MethodPost, backend.BackupRestore))
	http.HandleFunc("/backup/remove", Handle(http.MethodPost, backend.BackupRemove))
	http.HandleFunc("/installer/upgrade", Handle(http.MethodPost, backend.InstallerUpgrade))
	http.HandleFunc("/storage/disk_format", Handle(http.MethodPost, backend.StorageFormat))
	http.HandleFunc("/storage/boot_extend", Handle(http.MethodPost, backend.StorageBootExtend))
	http.HandleFunc("/event/trigger", Handle(http.MethodPost, backend.EventTrigger))
	http.HandleFunc("/redirect/domain/availability", Handle(http.MethodPost, backend.RedirectCheckFreeDomain))

	server := http.Server{}

	log.Println("Started backend")
	_ = server.Serve(unixListener)

}

func fail(w http.ResponseWriter, err error) {
	appError := err.Error()
	response := model.Response{
		Success: false,
		Message: &appError,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		_, _ = fmt.Fprintf(w, string(responseJson))
	}
}

func success(w http.ResponseWriter, data interface{}) {
	response := model.Response{
		Success: true,
		Data:    &data,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		fail(w, err)
	} else {
		_, _ = fmt.Fprintf(w, string(responseJson))
	}
}

func Handle(method string, f func(w http.ResponseWriter, req *http.Request) (interface{}, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("request: %s\n", req.URL.Path)
		if req.Method != method {
			fail(w, errors.New(fmt.Sprintf("wrong method %s, should be %s", req.Method, method)))
		}
		w.Header().Add("Content-Type", "application/json")
		data, err := f(w, req)
		if err != nil {
			fail(w, err)
		} else {
			success(w, data)
		}
	}
}

func (backend *Backend) BackupList(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return backend.backup.List()
}

func (backend *Backend) BackupRemove(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	var request model.BackupRemoveRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("file is missing")
	}
	err = backend.backup.Remove(request.File)
	if err != nil {
		return nil, err
	}
	return "removed", nil
}

func (backend *Backend) BackupCreate(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	var request model.BackupCreateRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("app is missing")
	}
	_ = backend.Master.Offer(job.JobBackupCreate{App: request.App})
	return "submitted", nil
}

func (backend *Backend) BackupRestore(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	var request model.BackupRestoreRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("file is missing")
	}
	_ = backend.Master.Offer(job.JobBackupRestore{File: request.File})
	return "submitted", nil
}

func (backend *Backend) InstallerUpgrade(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	_ = backend.Master.Offer(job.JobInstallerUpgrade{})
	return "submitted", nil
}

func (backend *Backend) JobStatus(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return backend.Master.Status().String(), nil
}

func (backend *Backend) StorageFormat(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	var request model.StorageFormatRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("device is missing")
	}
	_ = backend.Master.Offer(job.JobStorageFormat{Device: request.Device})
	return "submitted", nil
}

func (backend *Backend) EventTrigger(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	var request model.EventTriggerRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("event is missing")
	}
	return "ok", backend.eventTrigger.RunEventOnAllApps(request.Event)
}

func (backend *Backend) RedirectCheckFreeDomain(_ http.ResponseWriter, req *http.Request) (interface{}, error) {
	var request model.RedirectCheckFreeDomainRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("cannot parse request")
	}
	return "OK", backend.redirect.DomainAvailability(request)
}

func (backend *Backend) StorageBootExtend(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	_ = backend.Master.Offer(job.JobStorageBootExtend{})
	return "submitted", nil
}
