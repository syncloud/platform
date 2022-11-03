package ioc

import (
	"github.com/golobby/container/v3"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/syncloud/platform/access"
	"github.com/syncloud/platform/activation"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/cert"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/cron"
	"github.com/syncloud/platform/date"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/info"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/network"
	"github.com/syncloud/platform/nginx"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/snap"
	"github.com/syncloud/platform/storage"
	"github.com/syncloud/platform/storage/btrfs"
	"github.com/syncloud/platform/systemd"
	"github.com/syncloud/platform/version"
	"go.uber.org/zap"
	"time"
)

const (
	CertificateLogger = "CertificateLogger"
)

func Init(userConfig string, systemConfig string, backupDir string) {
	logger := log.Default()

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
	Singleton(func() *zap.Logger { return logger })
	NamedSingleton(CertificateLogger, func() *zap.Logger {
		return logger.With(zap.String(log.CategoryKey, log.CategoryCertificate))
	})
	Singleton(func() *network.Interface { return network.New() })
	Singleton(func() *retryablehttp.Client { return retryablehttp.NewClient() })
	Singleton(func() *version.PlatformVersion { return version.New() })
	Singleton(func() *identification.Parser { return identification.New() })
	Singleton(func(logger *zap.Logger) *cli.Executor { return cli.NewExecutor(logger) })
	Singleton(func(logger *zap.Logger) *auth.SystemPasswordChanger { return auth.NewSystemPassword(logger) })

	Singleton(func(executor *cli.Executor) *snap.Service { return snap.NewService(executor) })
	Singleton(func(executor *cli.Executor, systemConfig *config.SystemConfig) *systemd.Control {
		return systemd.New(executor, systemConfig, logger)
	})
	Singleton(func(snapService *snap.Service, systemConfig *config.SystemConfig, userConfig *config.UserConfig, control *systemd.Control) *nginx.Nginx {
		return nginx.New(control, systemConfig, userConfig)
	})
	Singleton(func(userConfig *config.UserConfig) *info.Device {
		return info.New(userConfig)
	})
	Singleton(func(userConfig *config.UserConfig, identification *identification.Parser, iface *network.Interface,
		client *retryablehttp.Client, version *version.PlatformVersion) *redirect.Service {
		var certLogger *zap.Logger
		NamedResolve(&certLogger, CertificateLogger)
		return redirect.New(userConfig, identification, iface, client, version, logger)
	})
	Singleton(func() *date.RealProvider { return date.New() })
	Singleton(func(redirectService *redirect.Service, userConfig *config.UserConfig, systemConfig *config.SystemConfig) *cert.Certbot {
		var certLogger *zap.Logger
		NamedResolve(&certLogger, CertificateLogger)
		return cert.NewCertbot(redirectService, userConfig, systemConfig, certLogger)
	})
	Singleton(func(systemConfig *config.SystemConfig, provider *date.RealProvider) *cert.Fake {
		var certLogger *zap.Logger
		NamedResolve(&certLogger, CertificateLogger)
		return cert.NewFake(systemConfig, provider, cert.SubjectOrganization, cert.DefaultDuration, certLogger)
	})
	Singleton(func(systemConfig *config.SystemConfig, userConfig *config.UserConfig, provider *date.RealProvider, certbot *cert.Certbot, fakeCert *cert.Fake, nginxService *nginx.Nginx) *cert.CertificateGenerator {
		var certLogger *zap.Logger
		NamedResolve(&certLogger, CertificateLogger)
		return cert.New(systemConfig, userConfig, provider, certbot, fakeCert, nginxService, certLogger)
	})
	Singleton(func(certGenerator *cert.CertificateGenerator) *cron.CertificateJob {
		return cron.NewCertificateJob(certGenerator)
	})
	Singleton(func() snap.SnapdClient { return snap.NewClient() })
	Singleton(func(snapClient snap.SnapdClient, deviceInfo *info.Device, systemConfig *config.SystemConfig, client *retryablehttp.Client) *snap.Snapd {
		return snap.New(snapClient, deviceInfo, systemConfig, client, logger)
	})

	Singleton(func(snapd *snap.Snapd, executor *cli.Executor) *event.Trigger {
		return event.New(snapd, executor)
	})

	Singleton(func(userConfig *config.UserConfig, redirectService *redirect.Service, eventTrigger *event.Trigger, client *retryablehttp.Client, netInfo *network.Interface, logger *zap.Logger) *access.PortProbe {
		return access.NewProbe(userConfig, client, logger)
	})

	Singleton(func(probe *access.PortProbe, userConfig *config.UserConfig, redirectService *redirect.Service, eventTrigger *event.Trigger, netInfo *network.Interface, logger *zap.Logger) *access.ExternalAddress {
		return access.New(probe, userConfig, redirectService, eventTrigger, netInfo, logger)
	})

	Singleton(func(job *access.ExternalAddress) *cron.ExternalAddressJob {
		return cron.NewExternalAddressJob(job)
	})
	Singleton(func(job1 *cron.CertificateJob, job2 *cron.ExternalAddressJob, userConfig *config.UserConfig) *cron.Cron {
		return cron.New([]cron.Job{job1, job2}, time.Minute*5, userConfig)
	})
	Singleton(func() *job.SingleJobMaster { return job.NewMaster() })
	Singleton(func(master *job.SingleJobMaster, logger *zap.Logger) *job.Worker {
		return job.NewWorker(master, logger)
	})
	Singleton(func(logger *zap.Logger) *backup.Backup { return backup.New(backupDir, logger) })
	Singleton(func() *installer.Installer { return installer.New() })
	Singleton(func() *storage.Storage { return storage.New() })
	Singleton(func(snapService *snap.Service, systemConfig *config.SystemConfig, executor *cli.Executor, passwordChanger *auth.SystemPasswordChanger) *auth.Service {
		return auth.New(snapService, systemConfig.DataDir(), systemConfig.AppDir(), systemConfig.ConfigDir(), executor, passwordChanger)
	})
	Singleton(func(ldapService *auth.Service, nginxService *nginx.Nginx, userConfig *config.UserConfig, eventTrigger *event.Trigger) *activation.Device {
		return activation.NewDevice(userConfig, ldapService, nginxService, eventTrigger)
	})
	Singleton(func() connection.InternetChecker { return connection.NewInternetChecker() })
	Singleton(func(internetChecker connection.InternetChecker, userConfig *config.UserConfig,
		redirectService *redirect.Service, device *activation.Device, certGenerator *cert.CertificateGenerator,
		logger *zap.Logger,
	) *activation.Managed {
		return activation.NewManaged(internetChecker, userConfig, redirectService, device, certGenerator, logger)
	})
	Singleton(func(internetChecker connection.InternetChecker, userConfig *config.UserConfig, device *activation.Device,
		certGenerator *cert.CertificateGenerator, logger *zap.Logger) *activation.Custom {
		return activation.NewCustom(internetChecker, userConfig, device, certGenerator, logger)
	})
	Singleton(func(activationManaged *activation.Managed, activationCustom *activation.Custom) *rest.Activate {
		return rest.NewActivateBackend(activationManaged, activationCustom)
	})
	Singleton(func(executor *cli.Executor) *systemd.JournalCtl { return systemd.NewJournalCtl(executor) })

	Singleton(func(certGenerator *cert.CertificateGenerator, journalCtl *systemd.JournalCtl) *rest.Certificate {
		return rest.NewCertificate(certGenerator, journalCtl)
	})

	Singleton(func(config *config.SystemConfig) *storage.PathChecker { return storage.NewPathChecker(config, logger) })
	Singleton(func(systemConfig *config.SystemConfig, executor *cli.Executor, checker *storage.PathChecker) *storage.Lsblk {
		return storage.NewLsblk(systemConfig, checker, executor, logger)
	})
	Singleton(func(executor *cli.Executor) *storage.FreeSpaceChecker { return storage.NewFreeSpaceChecker(executor) })
	Singleton(func() *storage.Linker { return storage.NewLinker() })
	Singleton(func(systemConfig *config.SystemConfig, executor *cli.Executor) *btrfs.Stats {
		return btrfs.NewStats(systemConfig, executor)
	})
	Singleton(func(systemConfig *config.SystemConfig, executor *cli.Executor, stats *btrfs.Stats, systemd *systemd.Control) *btrfs.Disks {
		return btrfs.NewDisks(systemConfig, executor, systemd, logger)
	})
	Singleton(func(systemConfig *config.SystemConfig, freeSpaceChecker *storage.FreeSpaceChecker,
		systemd *systemd.Control, eventTrigger *event.Trigger, lsblk *storage.Lsblk,
		executor *cli.Executor, linker *storage.Linker, btrfs *btrfs.Disks, stats *btrfs.Stats) *storage.Disks {
		return storage.NewDisks(systemConfig, eventTrigger, lsblk, systemd, freeSpaceChecker, linker, executor, btrfs, stats, logger)
	})

	Singleton(func(master *job.SingleJobMaster, backupService *backup.Backup, eventTrigger *event.Trigger, worker *job.Worker,
		redirectService *redirect.Service, installerService *installer.Installer, storageService *storage.Storage,
		id *identification.Parser, activate *rest.Activate, userConfig *config.UserConfig, cert *rest.Certificate,
		externalAddress *access.ExternalAddress, snapd *snap.Snapd, disks *storage.Disks, journalCtl *systemd.JournalCtl,
	) *rest.Backend {
		return rest.NewBackend(master, backupService, eventTrigger, worker, redirectService,
			installerService, storageService, id, activate, userConfig, cert, externalAddress, snapd, disks, journalCtl)
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
