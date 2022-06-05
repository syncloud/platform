package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest/model"
	"github.com/syncloud/platform/snap"
	"github.com/syncloud/platform/storage"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/syncloud/platform/access"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/job"
)

type Backend struct {
	JobMaster       *job.Master
	backup          *backup.Backup
	eventTrigger    *event.Trigger
	worker          *job.Worker
	redirect        *redirect.Service
	installer       installer.AppInstaller
	storage         *storage.Storage
	identification  *identification.Parser
	activate        *Activate
	userConfig      *config.UserConfig
	certificate     *Certificate
	externalAddress *access.ExternalAddress
	snapd           *snap.Snapd
	disks           *storage.Disks
}

func NewBackend(master *job.Master, backup *backup.Backup,
	eventTrigger *event.Trigger, worker *job.Worker,
	redirect *redirect.Service, installerService *installer.Installer,
	storageService *storage.Storage,
	identification *identification.Parser,
	activate *Activate, userConfig *config.UserConfig,
	certificate *Certificate, externalAddress *access.ExternalAddress,
	snapd *snap.Snapd, disks *storage.Disks,
) *Backend {

	return &Backend{
		JobMaster:       master,
		backup:          backup,
		eventTrigger:    eventTrigger,
		worker:          worker,
		redirect:        redirect,
		installer:       installerService,
		storage:         storageService,
		identification:  identification,
		activate:        activate,
		userConfig:      userConfig,
		certificate:     certificate,
		externalAddress: externalAddress,
		snapd:           snapd,
		disks:           disks,
	}
}

func (b *Backend) NewReverseProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		redirectApiUrl := b.userConfig.GetRedirectApiUrl()
		redirectUrl, err := url.Parse(redirectApiUrl)
		if err != nil {
			fmt.Printf("proxy url error: %v", err)
			return
		}
		fmt.Printf("proxy url: %v", redirectUrl)

		req.URL.Scheme = redirectUrl.Scheme
		req.URL.Host = redirectUrl.Host
		req.Host = redirectUrl.Host
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
	r.HandleFunc("/installer/version", Handle(b.InstallerVersion)).Methods("GET")
	r.HandleFunc("/storage/disk_format", Handle(b.StorageFormat)).Methods("POST")
	r.HandleFunc("/storage/boot_extend", Handle(b.StorageBootExtend)).Methods("POST")
	r.HandleFunc("/storage/boot/disk", Handle(b.StorageBootDisk)).Methods("GET")
	r.HandleFunc("/storage/disk/deactivate", Handle(b.StorageDiskDeactivate)).Methods("POST")
	r.HandleFunc("/storage/disks", Handle(b.StorageDisks)).Methods("GET")
	r.HandleFunc("/event/trigger", Handle(b.EventTrigger)).Methods("POST")
	r.HandleFunc("/activate/managed", Handle(b.activate.Managed)).Methods("POST")
	r.HandleFunc("/activate/custom", Handle(b.activate.Custom)).Methods("POST")
	r.HandleFunc("/id", Handle(b.Id)).Methods("GET")
	r.HandleFunc("/certificate", Handle(b.certificate.Certificate)).Methods("GET")
	r.HandleFunc("/certificate/log", Handle(b.certificate.CertificateLog)).Methods("GET")
	r.HandleFunc("/redirect_info", Handle(b.RedirectInfo)).Methods("GET")
	r.HandleFunc("/access", Handle(b.GetAccess)).Methods("GET")
	r.HandleFunc("/access", Handle(b.SetAccess)).Methods("POST")
	r.HandleFunc("/activation/status", Handle(b.IsActivated)).Methods("GET")
	r.HandleFunc("/apps/available", Handle(b.AppsAvailable)).Methods("GET")
	r.HandleFunc("/apps/installed", Handle(b.AppsInstalled)).Methods("GET")
	r.PathPrefix("/redirect/domain/availability").Handler(http.StripPrefix("/redirect", b.NewReverseProxy()))
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	r.Use(middleware)

	fmt.Println("Started backend")
	_ = http.Serve(listener, r)

}

func fail(w http.ResponseWriter, err error) {
	fmt.Println("error: ", err)
	response := model.Response{
		Success: false,
		Message: err.Error(),
	}
	statusCode := http.StatusInternalServerError
	switch v := err.(type) {
	case *model.ParameterError:
		fmt.Println("parameter error: ", v.ParameterErrors)
		response.ParametersMessages = v.ParameterErrors
		statusCode = 400
	}
	responseJson, err := json.Marshal(response)
	responseText := ""
	if err != nil {
		responseText = err.Error()
	} else {
		responseText = string(responseJson)
	}
	http.Error(w, responseText, statusCode)
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
		fmt.Printf("%s: %s\n", r.Method, r.RequestURI)
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("404 %s: %s\n", r.Method, r.RequestURI)
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
		fmt.Printf("parse error: %v\n", err.Error())
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
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("app is missing")
	}
	_ = b.JobMaster.Offer(func() { b.backup.Create(request.App) })
	return "submitted", nil
}

func (b *Backend) BackupRestore(req *http.Request) (interface{}, error) {
	var request model.BackupRestoreRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("file is missing")
	}
	_ = b.JobMaster.Offer(func() { b.backup.Restore(request.File) })
	return "submitted", nil
}

func (b *Backend) InstallerUpgrade(_ *http.Request) (interface{}, error) {
	_ = b.JobMaster.Offer(func() { b.installer.Upgrade() })
	return "submitted", nil
}

func (b *Backend) JobStatus(_ *http.Request) (interface{}, error) {
	return b.JobMaster.Status().String(), nil
}

func (b *Backend) StorageFormat(req *http.Request) (interface{}, error) {
	var request model.StorageFormatRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("device is missing")
	}
	_ = b.JobMaster.Offer(func() { b.storage.Format(request.Device) })
	return "submitted", nil
}

func (b *Backend) EventTrigger(req *http.Request) (interface{}, error) {
	var request model.EventTriggerRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("event is missing")
	}
	return "ok", b.eventTrigger.RunEventOnAllApps(request.Event)
}

func (b *Backend) RedirectInfo(_ *http.Request) (interface{}, error) {
	fmt.Printf("redirect info\n")
	response := &model.RedirectInfoResponse{
		Domain: b.userConfig.GetRedirectDomain(),
	}
	return response, nil
}

func (b *Backend) GetAccess(_ *http.Request) (interface{}, error) {
	response := &model.Access{
		Ipv4:        b.userConfig.GetPublicIp(),
		Ipv4Enabled: b.userConfig.IsIpv4Enabled(),
		Ipv4Public:  b.userConfig.IsIpv4Public(),
		AccessPort:  b.userConfig.GetPublicPort(),
		Ipv6Enabled: b.userConfig.IsIpv6Enabled(),
	}
	return response, nil
}

func (b *Backend) IsActivated(_ *http.Request) (interface{}, error) {
	return b.userConfig.IsActivated(), nil
}

func (b *Backend) SetAccess(req *http.Request) (interface{}, error) {
	var request model.Access
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("access request is wrong")
	}

	return request, b.externalAddress.Update(request)
}

func (b *Backend) AppsAvailable(_ *http.Request) (interface{}, error) {
	return b.snapd.StoreUserApps()
}

func (b *Backend) AppsInstalled(_ *http.Request) (interface{}, error) {
	return b.snapd.InstalledUserApps()
}

func (b *Backend) InstallerVersion(_ *http.Request) (interface{}, error) {
	return b.snapd.Installer()
}

func (b *Backend) Id(_ *http.Request) (interface{}, error) {
	id, err := b.identification.Id()
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("id is not available")
	}
	return id, nil
}

func (b *Backend) StorageBootExtend(_ *http.Request) (interface{}, error) {
	_ = b.JobMaster.Offer(func() { b.storage.BootExtend() })
	return "submitted", nil
}

func (b *Backend) StorageDisks(_ *http.Request) (interface{}, error) {
	return b.disks.AvailableDisks()
}

func (b *Backend) StorageBootDisk(_ *http.Request) (interface{}, error) {
	return b.disks.RootPartition()
}

func (b *Backend) StorageDiskDeactivate(_ *http.Request) (interface{}, error) {
	return "OK", b.disks.DeactivateDisk()
}
