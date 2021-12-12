package ioc

import (
	"github.com/golobby/container/v3"
	"github.com/syncloud/platform/activation"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/cert"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/cron"
	"github.com/syncloud/platform/date"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/network"
	"github.com/syncloud/platform/nginx"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/snap"
	"github.com/syncloud/platform/storage"
	"github.com/syncloud/platform/systemd"
	"go.uber.org/zap"
	"time"
)

const (
	CertificateLogger = "CertificateLogger"
)

func Init(userConfig string, systemConfig string, backupDir string, logger *zap.Logger) {

	Singleton(func() *config.UserConfig {
		userConfig := config.NewUserConfig(userConfig, config.OldConfig)
		userConfig.Load()
		return userConfig
	})
	Singleton(func() *config.SystemConfig {
		systemConfig := config.NewSystemConfig(systemConfig)
		systemConfig.Load()
		return systemConfig
	})
	Singleton(func() *network.Interface { return network.New() })
	Singleton(func() *identification.Parser { return identification.New() })
	Singleton(func(userConfig *config.UserConfig, identification *identification.Parser, iface *network.Interface) *redirect.Service {
		return redirect.New(userConfig, identification, iface)
	})
	Singleton(func() *date.RealProvider { return date.New() })
	Singleton(func(redirectService *redirect.Service, userConfig *config.UserConfig, systemConfig *config.SystemConfig) *cert.Certbot {
		return cert.NewCertbot(redirectService, userConfig, systemConfig)
	})
	NamedSingleton(CertificateLogger, func() *zap.Logger {
		return logger.With(zap.String("category", "certificate"))
	})
	Singleton(func(systemConfig *config.SystemConfig) *cert.Fake {
		var certLogger *zap.Logger
		NamedResolve(&certLogger, CertificateLogger)
		return cert.NewFake(systemConfig, certLogger)
	})
	Singleton(func(systemConfig *config.SystemConfig, userConfig *config.UserConfig, provider *date.RealProvider, certbot *cert.Certbot, fakeCert *cert.Fake) *cert.CertificateGenerator {
		var certLogger *zap.Logger
		NamedResolve(&certLogger, CertificateLogger)
		return cert.New(systemConfig, userConfig, provider, certbot, fakeCert, certLogger)
	})
	Singleton(func(certGenerator *cert.CertificateGenerator) *cron.CertificateJob {
		return cron.NewCertificateJob(certGenerator)
	})
	Singleton(func(userConfig *config.UserConfig) *cron.PortsJob {
		return cron.NewPortsJob(userConfig)
	})
	Singleton(func(job1 *cron.CertificateJob, job2 *cron.PortsJob, userConfig *config.UserConfig) *cron.Cron {
		return cron.New([]cron.Job{job1, job2}, time.Minute*5, userConfig)
	})
	Singleton(func() *job.Master { return job.NewMaster() })
	Singleton(func(master *job.Master) *job.Worker { return job.NewWorker(master) })
	Singleton(func() *backup.Backup { return backup.New(backupDir) })
	Singleton(func() snap.SnapdClient { return snap.NewClient() })
	Singleton(func(snapClient snap.SnapdClient) *snap.Snapd { return snap.New(snapClient) })
	Singleton(func(snapd *snap.Snapd) *event.Trigger { return event.New(snapd) })
	Singleton(func() *installer.Installer { return installer.New() })
	Singleton(func() *storage.Storage { return storage.New() })
	Singleton(func() *snap.Service { return snap.NewService() })
	Singleton(func(snapService *snap.Service, systemConfig *config.SystemConfig) *auth.Service {
		return auth.New(snapService, systemConfig.DataDir(), systemConfig.AppDir(), systemConfig.ConfigDir())
	})
	Singleton(func(snapService *snap.Service, systemConfig *config.SystemConfig, userConfig *config.UserConfig) *nginx.Nginx {
		return nginx.New(systemd.New(), systemConfig, userConfig)
	})
	Singleton(func(ldapService *auth.Service, nginxService *nginx.Nginx, userConfig *config.UserConfig, eventTrigger *event.Trigger) *activation.Device {
		return activation.NewDevice(userConfig, ldapService, nginxService, eventTrigger)
	})
	Singleton(func() connection.InternetChecker { return connection.NewInternetChecker() })
	Singleton(func(internetChecker connection.InternetChecker, userConfig *config.UserConfig,
		redirectService *redirect.Service, device *activation.Device, certGenerator *cert.CertificateGenerator,
	) *activation.Managed {
		return activation.NewManaged(internetChecker, userConfig, redirectService, device, certGenerator)
	})
	Singleton(func(internetChecker connection.InternetChecker, userConfig *config.UserConfig, device *activation.Device,
		certGenerator *cert.CertificateGenerator) *activation.Custom {
		return activation.NewCustom(internetChecker, userConfig, device, certGenerator)
	})
	Singleton(func(activationManaged *activation.Managed, activationCustom *activation.Custom) *rest.Activate {
		return rest.NewActivateBackend(activationManaged, activationCustom)
	})
	Singleton(func() *cert.Reader { return cert.NewReader() })
	Singleton(func(master *job.Master, backupService *backup.Backup, eventTrigger *event.Trigger, worker *job.Worker,
		redirectService *redirect.Service, installerService *installer.Installer, storageService *storage.Storage,
		id *identification.Parser, activate *rest.Activate, userConfig *config.UserConfig, certReader *cert.Reader,
	) *rest.Backend {
		return rest.NewBackend(master, backupService, eventTrigger, worker, redirectService,
			installerService, storageService, id, activate, userConfig, certReader)
	})

}

func Singleton(resolver interface{}) {
	err := container.Singleton(resolver)
	if err != nil {
		panic(err)
	}
}

func NamedSingleton(name string, resolver interface{}) {
	err := container.NamedSingleton(name, resolver)
	if err != nil {
		panic(err)
	}
}

func Call(abstraction interface{}) {
	err := container.Call(abstraction)
	if err != nil {
		panic(err)
	}
}

func NamedResolve(abstraction interface{}, name string) {
	err := container.NamedResolve(abstraction, name)
	if err != nil {
		panic(err)
	}
}
