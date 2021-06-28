package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/syncloud/platform/activation"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest/model"
	"github.com/syncloud/platform/storage"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/job"
)

type Backend struct {
	Master           *job.Master
	backup           *backup.Backup
	eventTrigger     *event.Trigger
	worker           *job.Worker
	redirect         *redirect.Service
	installer        installer.AppInstaller
	storage          *storage.Storage
	redirectProxy    *httputil.ReverseProxy
	identification   *identification.Parser
	activationFree   *activation.Free
	activationCustom *activation.Custom
}

func NewBackend(master *job.Master, backup *backup.Backup,
	eventTrigger *event.Trigger, worker *job.Worker,
	redirect *redirect.Service, installerService *installer.Installer,
	storageService *storage.Storage, redirectUrl *url.URL,
	identification *identification.Parser,
	activationFree *activation.Free,
	activationCustom *activation.Custom,
) *Backend {

	return &Backend{
		Master:           master,
		backup:           backup,
		eventTrigger:     eventTrigger,
		worker:           worker,
		redirect:         redirect,
		installer:        installerService,
		storage:          storageService,
		redirectProxy:    NewReverseProxy(redirectUrl),
		identification:   identification,
		activationFree:   activationFree,
		activationCustom: activationCustom,
	}
}

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}
	return &httputil.ReverseProxy{Director: director}
}

func (b *Backend) Start(network string, address string) {
	listener, err := net.Listen(network, address)
	if err != nil {
		panic(err)
	}

	go b.worker.Start()

	r := mux.NewRouter()
	r.HandleFunc("/job/status", Handle(b.JobStatus)).Methods("GET")
	r.HandleFunc("/backup/list", Handle(b.BackupList)).Methods("GET")
	r.HandleFunc("/backup/create", Handle(b.BackupCreate)).Methods("POST")
	r.HandleFunc("/backup/restore", Handle(b.BackupRestore)).Methods("POST")
	r.HandleFunc("/backup/remove", Handle(b.BackupRemove)).Methods("POST")
	r.HandleFunc("/installer/upgrade", Handle(b.InstallerUpgrade)).Methods("POST")
	r.HandleFunc("/storage/disk_format", Handle(b.StorageFormat)).Methods("POST")
	r.HandleFunc("/storage/boot_extend", Handle(b.StorageBootExtend)).Methods("POST")
	r.HandleFunc("/event/trigger", Handle(b.EventTrigger)).Methods("POST")
	r.HandleFunc("/activate/free", Handle(b.ActivateFree)).Methods("POST")
	r.HandleFunc("/activate/custom", Handle(b.ActivateCustom)).Methods("POST")
	r.HandleFunc("/id", Handle(b.Id)).Methods("GET")
	r.PathPrefix("/redirect").Handler(http.StripPrefix("/redirect", b.redirectProxy))
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	r.Use(middleware)

	log.Println("Started backend")
	_ = http.Serve(listener, r)

}

func fail(w http.ResponseWriter, err error) {
	log.Println("error: ", err)
	response := model.Response{
		Success: false,
		Message: err.Error(),
	}
	statusCode := http.StatusInternalServerError
	switch v := err.(type) {
	case *model.ParameterError:
		response.ParametersMessages = v.ParameterErrors
		statusCode = 400
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
	} else {
		http.Error(w, string(responseJson), statusCode)
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

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.RequestURI)
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("404 %s: %s", r.Method, r.RequestURI)
	http.NotFound(w, r)
}

func Handle(f func(req *http.Request) (interface{}, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := f(req)
		if err != nil {
			fail(w, err)
		} else {
			success(w, data)
		}
	}
}

func (b *Backend) BackupList(_ *http.Request) (interface{}, error) {
	return b.backup.List()
}

func (b *Backend) BackupRemove(req *http.Request) (interface{}, error) {
	var request model.BackupRemoveRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("file is missing")
	}
	err = b.backup.Remove(request.File)
	if err != nil {
		return nil, err
	}
	return "removed", nil
}

func (b *Backend) BackupCreate(req *http.Request) (interface{}, error) {
	var request model.BackupCreateRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("app is missing")
	}
	_ = b.Master.Offer(func() { b.backup.Create(request.App) })
	return "submitted", nil
}

func (b *Backend) BackupRestore(req *http.Request) (interface{}, error) {
	var request model.BackupRestoreRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("file is missing")
	}
	_ = b.Master.Offer(func() { b.backup.Restore(request.File) })
	return "submitted", nil
}

func (b *Backend) InstallerUpgrade(_ *http.Request) (interface{}, error) {
	_ = b.Master.Offer(func() { b.installer.Upgrade() })
	return "submitted", nil
}

func (b *Backend) JobStatus(_ *http.Request) (interface{}, error) {
	return b.Master.Status().String(), nil
}

func (b *Backend) StorageFormat(req *http.Request) (interface{}, error) {
	var request model.StorageFormatRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("device is missing")
	}
	_ = b.Master.Offer(func() { b.storage.Format(request.Device) })
	return "submitted", nil
}

func (b *Backend) EventTrigger(req *http.Request) (interface{}, error) {
	var request model.EventTriggerRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("event is missing")
	}
	return "ok", b.eventTrigger.RunEventOnAllApps(request.Event)
}

func (b *Backend) ActivateFree(req *http.Request) (interface{}, error) {
	var request activation.FreeActivateRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return "ok", b.activationFree.Activate(request.RedirectEmail, request.RedirectPassword, request.Domain, request.DeviceUsername, request.DevicePassword)
}

func (b *Backend) ActivateCustom(req *http.Request) (interface{}, error) {
	var request activation.CustomActivateRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	if len(request.DeviceUsername) < 3 {
		return nil, model.SingleParameterError("device_username", "less than 3 characters")
	}
	if len(request.DevicePassword) < 7 {
		return nil, model.SingleParameterError("device_password", "less than 7 characters")
	}

	return "ok", b.activationCustom.Activate(request.Domain, request.DeviceUsername, request.DevicePassword)
}

func (b *Backend) Id(_ *http.Request) (interface{}, error) {
	id, err := b.identification.Id()
	if err != nil {
		log.Printf("parse error: %v", err.Error())
		return nil, errors.New("id is not available")
	}
	return id, nil
}

func (b *Backend) StorageBootExtend(_ *http.Request) (interface{}, error) {
	_ = b.Master.Offer(func() { b.storage.BootExtend() })
	return "submitted", nil
}
