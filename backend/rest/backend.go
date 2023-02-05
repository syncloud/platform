package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/info"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/network"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest/model"
	"github.com/syncloud/platform/snap"
	"github.com/syncloud/platform/storage"
	"github.com/syncloud/platform/support"
	"github.com/syncloud/platform/systemd"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/syncloud/platform/access"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/job"
)

type Backend struct {
	JobMaster       *job.SingleJobMaster
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
	snapd           *snap.Server
	disks           *storage.Disks
	journalCtl      *systemd.Journal
	deviceInfo      *info.Device
	executor        *cli.ShellExecutor
	iface           *network.TcpInterfaces
	support         *support.Sender
}

func NewBackend(
	master *job.SingleJobMaster, backup *backup.Backup, eventTrigger *event.Trigger, worker *job.Worker,
	redirect *redirect.Service, installerService *installer.Installer, storageService *storage.Storage,
	identification *identification.Parser, activate *Activate, userConfig *config.UserConfig,
	certificate *Certificate, externalAddress *access.ExternalAddress, snapd *snap.Server,
	disks *storage.Disks, journalCtl *systemd.Journal, deviceInfo *info.Device, executor *cli.ShellExecutor,
	iface *network.TcpInterfaces, support *support.Sender) *Backend {

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
		journalCtl:      journalCtl,
		deviceInfo:      deviceInfo,
		executor:        executor,
		iface:           iface,
		support:         support,
	}
}

func (b *Backend) Start(network string, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	go b.worker.Start()

	r := mux.NewRouter()
	r.HandleFunc("/job/status", Handle(b.JobStatus)).Methods("GET")
	r.HandleFunc("/backup/list", Handle(b.BackupList)).Methods("GET")
	r.HandleFunc("/backup/auto", Handle(b.GetBackupAuto)).Methods("GET")
	r.HandleFunc("/backup/auto", Handle(b.SetBackupAuto)).Methods("POST")
	r.HandleFunc("/backup/create", Handle(b.BackupCreate)).Methods("POST")
	r.HandleFunc("/backup/restore", Handle(b.BackupRestore)).Methods("POST")
	r.HandleFunc("/backup/remove", Handle(b.BackupRemove)).Methods("POST")
	r.HandleFunc("/installer/upgrade", Handle(b.InstallerUpgrade)).Methods("POST")
	r.HandleFunc("/installer/version", Handle(b.InstallerVersion)).Methods("GET")
	r.HandleFunc("/installer/status", Handle(b.InstallerStatus)).Methods("GET")
	r.HandleFunc("/storage/boot_extend", Handle(b.StorageBootExtend)).Methods("POST")
	r.HandleFunc("/storage/boot/disk", Handle(b.StorageBootDisk)).Methods("GET")
	r.HandleFunc("/storage/deactivate", Handle(b.StorageDiskDeactivate)).Methods("POST")
	r.HandleFunc("/storage/activate/partition", Handle(b.StorageActivatePartition)).Methods("POST")
	r.HandleFunc("/storage/activate/disk", Handle(b.StorageActivateDisks)).Methods("POST")
	r.HandleFunc("/storage/error/last", Handle(b.StorageLastError)).Methods("GET")
	r.HandleFunc("/storage/error/clear", Handle(b.StorageClearError)).Methods("POST")
	r.HandleFunc("/storage/disks", Handle(b.StorageDisks)).Methods("GET")
	r.HandleFunc("/event/trigger", Handle(b.EventTrigger)).Methods("POST")
	r.HandleFunc("/activate/managed", Handle(b.activate.Managed)).Methods("POST")
	r.HandleFunc("/activate/custom", Handle(b.activate.Custom)).Methods("POST")
	r.HandleFunc("/deactivate", Handle(b.Deactivate)).Methods("POST")
	r.HandleFunc("/id", Handle(b.Id)).Methods("GET")
	r.HandleFunc("/certificate", Handle(b.certificate.Certificate)).Methods("GET")
	r.HandleFunc("/certificate/log", Handle(b.certificate.CertificateLog)).Methods("GET")
	r.HandleFunc("/redirect_info", Handle(b.RedirectInfo)).Methods("GET")
	r.HandleFunc("/access", Handle(b.GetAccess)).Methods("GET")
	r.HandleFunc("/access", Handle(b.SetAccess)).Methods("POST")
	r.HandleFunc("/activation/status", Handle(b.IsActivated)).Methods("GET")
	r.HandleFunc("/apps/available", Handle(b.AppsAvailable)).Methods("GET")
	r.HandleFunc("/apps/installed", Handle(b.AppsInstalled)).Methods("GET")
	r.HandleFunc("/app/install", Handle(b.AppInstall)).Methods("POST")
	r.HandleFunc("/app/remove", Handle(b.AppRemove)).Methods("POST")
	r.HandleFunc("/app/upgrade", Handle(b.AppUpgrade)).Methods("POST")
	r.HandleFunc("/app", Handle(b.App)).Methods("GET")
	r.HandleFunc("/logs", Handle(b.Logs)).Methods("GET")
	r.HandleFunc("/logs/send", Handle(b.SendLogs)).Methods("POST")
	r.HandleFunc("/device/url", Handle(b.DeviceUrl)).Methods("GET")
	r.HandleFunc("/restart", Handle(b.Restart)).Methods("POST")
	r.HandleFunc("/shutdown", Handle(b.Shutdown)).Methods("POST")
	r.HandleFunc("/network/interfaces", Handle(b.NetworkInterfaces)).Methods("GET")

	redirectProxy, err := b.ProxyRedirect()
	if err != nil {
		return err
	}
	r.PathPrefix("/redirect/domain/availability").Handler(http.StripPrefix("/redirect", redirectProxy))

	r.PathPrefix("/proxy/image").Handler(b.ProxyImage())
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	r.Use(middleware)

	fmt.Println("Started backend")
	_ = http.Serve(listener, r)
	return nil
}

func (b *Backend) ProxyRedirect() (*httputil.ReverseProxy, error) {
	redirectApiUrl := b.userConfig.GetRedirectApiUrl()
	redirectUrl, err := url.Parse(redirectApiUrl)
	if err != nil {
		fmt.Printf("proxy url error: %v", err)
		return nil, err
	}
	fmt.Printf("proxy url: %v", redirectUrl)
	director := func(req *http.Request) {
		req.URL.Scheme = redirectUrl.Scheme
		req.URL.Host = redirectUrl.Host
		req.Host = redirectUrl.Host
	}
	return &httputil.ReverseProxy{Director: director}, nil
}

func (b *Backend) ProxyImage() *httputil.ReverseProxy {
	host := "apps.syncloud.org"
	director := func(req *http.Request) {
		query := req.URL.Query()
		if !query.Has("channel") {
			return
		}
		if !query.Has("app") {
			return
		}
		req.URL.Scheme = "http"
		req.URL.Host = host
		req.URL.Path = fmt.Sprintf("releases/%s/images/%s-128.png", query.Get("channel"), query.Get("app"))
		req.Host = host
	}
	return &httputil.ReverseProxy{Director: director}
}

func (b *Backend) BackupList(_ *http.Request) (interface{}, error) {
	return b.backup.List()
}

func (b *Backend) GetBackupAuto(_ *http.Request) (interface{}, error) {
	return b.backup.Auto(), nil
}

func (b *Backend) SetBackupAuto(req *http.Request) (interface{}, error) {
	var request backup.Auto
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("bad request")
	}
	b.backup.SetAuto(request)
	return "OK", nil
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
	_ = b.JobMaster.Offer("backup.create", func() error { return b.backup.Create(request.App) })
	return "submitted", nil
}

func (b *Backend) BackupRestore(req *http.Request) (interface{}, error) {
	var request model.BackupRestoreRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("file is missing")
	}
	_ = b.JobMaster.Offer("backup.restore", func() error { return b.backup.Restore(request.File) })
	return "submitted", nil
}

func (b *Backend) InstallerUpgrade(_ *http.Request) (interface{}, error) {
	_ = b.JobMaster.Offer("installer.upgrade", func() error { return b.installer.Upgrade() })
	return "submitted", nil
}

func (b *Backend) JobStatus(_ *http.Request) (interface{}, error) {
	return b.JobMaster.Status(), nil
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

func (b *Backend) AppUpgrade(req *http.Request) (interface{}, error) {
	var request model.AppActionRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("wrong request")
	}

	return nil, b.snapd.Upgrade(request.AppId)
}

func (b *Backend) AppInstall(req *http.Request) (interface{}, error) {
	var request model.AppActionRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("wrong request")
	}

	return nil, b.snapd.Install(request.AppId)
}

func (b *Backend) AppRemove(req *http.Request) (interface{}, error) {
	var request model.AppActionRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("wrong request")
	}

	return nil, b.snapd.Remove(request.AppId)
}

func (b *Backend) App(req *http.Request) (interface{}, error) {
	query := req.URL.Query()
	if query.Has("app_id") {
		id := query.Get("app_id")
		return b.snapd.Find(id)
	} else {
		return nil, errors.New("app_id query param is missing")
	}
}

func (b *Backend) InstallerVersion(_ *http.Request) (interface{}, error) {
	return b.snapd.Installer()
}

func (b *Backend) InstallerStatus(_ *http.Request) (interface{}, error) {
	return b.snapd.Changes()
}

func (b *Backend) Id(_ *http.Request) (interface{}, error) {
	id, err := b.identification.Id()
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("id is not available")
	}
	return id, nil
}

func (b *Backend) Deactivate(_ *http.Request) (interface{}, error) {
	b.userConfig.SetDeactivated()
	return "OK", nil
}

func (b *Backend) StorageBootExtend(_ *http.Request) (interface{}, error) {
	_ = b.JobMaster.Offer("storage.boot.extend", func() error { return b.storage.BootExtend() })
	return "submitted", nil
}

func (b *Backend) StorageDisks(_ *http.Request) (interface{}, error) {
	return b.disks.AvailableDisks()
}

func (b *Backend) StorageBootDisk(_ *http.Request) (interface{}, error) {
	return b.disks.RootPartition()
}

func (b *Backend) StorageDiskDeactivate(_ *http.Request) (interface{}, error) {
	return "OK", b.disks.Deactivate()
}

func (b *Backend) StorageActivatePartition(req *http.Request) (interface{}, error) {
	var request model.StorageActivatePartitionRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, err
	}
	if request.Format {
		err = b.storage.Format(request.Device)
		if err != nil {
			fmt.Printf("format error: %v\n", err.Error())
			return nil, err
		}
	}

	return "OK", b.JobMaster.Offer("storage.activate.partition", func() error { return b.disks.ActivatePartition(request.Device) })

}

func (b *Backend) StorageActivateDisks(req *http.Request) (interface{}, error) {
	var request model.StorageActivateDisksRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, err
	}

	return "OK", b.JobMaster.Offer("storage.activate.disks", func() error { return b.disks.ActivateDisks(request.Devices, request.Format) })
}

func (b *Backend) StorageLastError(_ *http.Request) (interface{}, error) {
	return "OK", b.disks.GetLastError()
}

func (b *Backend) StorageClearError(_ *http.Request) (interface{}, error) {
	b.disks.ClearLastError()
	return "OK", nil
}

func (b *Backend) Logs(_ *http.Request) (interface{}, error) {
	return b.journalCtl.ReadAll(func(line string) bool {
		return true
	}), nil
}

func (b *Backend) DeviceUrl(_ *http.Request) (interface{}, error) {
	return b.deviceInfo.DeviceUrl(), nil
}

func (b *Backend) Restart(_ *http.Request) (interface{}, error) {
	return b.executor.CombinedOutput("shutdown", "-r", "now")
}

func (b *Backend) Shutdown(_ *http.Request) (interface{}, error) {
	return b.executor.CombinedOutput("shutdown", "now")
}

func (b *Backend) NetworkInterfaces(_ *http.Request) (interface{}, error) {
	return b.iface.List()
}

func (b *Backend) SendLogs(req *http.Request) (interface{}, error) {
	includeSupport := false
	query := req.URL.Query()
	if query.Has("include_support") {
		includeSupport = query.Get("app_id") == "true"
	}
	return b.support.Send(includeSupport), nil
}
