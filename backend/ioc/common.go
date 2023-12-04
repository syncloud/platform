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
	"github.com/syncloud/platform/du"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/hook"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/info"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/network"
	"github.com/syncloud/platform/nginx"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/session"
	"github.com/syncloud/platform/snap"
	"github.com/syncloud/platform/storage"
	"github.com/syncloud/platform/storage/btrfs"
	"github.com/syncloud/platform/support"
	"github.com/syncloud/platform/systemd"
	"github.com/syncloud/platform/version"
	"go.uber.org/zap"
	"time"
)

const (
	CertificateLogger = "CertificateLogger"
)

var logger = log.Default()

func Init(userConfig string, systemConfig string, backupDir string, varDir string) (container.Container, error) {
	c := container.New()
	err := c.Singleton(func() *config.UserConfig {
		userConfig := config.NewUserConfig(userConfig, config.OldConfig)
		userConfig.Load()
		return userConfig
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *config.SystemConfig {
		systemConfig := config.NewSystemConfig(systemConfig)
		systemConfig.Load()
		return systemConfig
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *zap.Logger { return logger })
	if err != nil {
		return nil, err
	}
	err = c.NamedSingleton(CertificateLogger, func() *zap.Logger {
		return logger.With(zap.String(log.CategoryKey, log.CategoryCertificate))
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *network.TcpInterfaces { return network.New() })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *retryablehttp.Client {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = 10
		return retryClient
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *version.PlatformVersion { return version.New() })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *identification.Parser { return identification.New() })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(logger *zap.Logger) *cli.ShellExecutor { return cli.New(logger) })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(logger *zap.Logger) *auth.SystemPasswordChanger { return auth.NewSystemPassword(logger) })

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(executor *cli.ShellExecutor, logger *zap.Logger) *snap.Cli { return snap.NewCli(executor, logger) })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(executor *cli.ShellExecutor, systemConfig *config.SystemConfig) *systemd.Control {
		return systemd.New(executor, systemConfig, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(systemConfig *config.SystemConfig, userConfig *config.UserConfig, control *systemd.Control) *nginx.Nginx {
		return nginx.New(control, systemConfig, userConfig)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(userConfig *config.UserConfig) *info.Device {
		return info.New(userConfig)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(
		userConfig *config.UserConfig,
		identification *identification.Parser,
		iface *network.TcpInterfaces,
		client *retryablehttp.Client,
		version *version.PlatformVersion,
	) (*redirect.Service, error) {
		var certLogger *zap.Logger
		err := c.NamedResolve(&certLogger, CertificateLogger)
		if err != nil {
			return nil, err
		}
		return redirect.New(
			userConfig,
			identification,
			iface,
			client,
			version,
			logger,
		), nil
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *date.RealProvider { return date.New() })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(redirectService *redirect.Service, userConfig *config.UserConfig, systemConfig *config.SystemConfig) (*cert.Certbot, error) {
		var certLogger *zap.Logger
		err := c.NamedResolve(&certLogger, CertificateLogger)
		if err != nil {
			return nil, err
		}
		return cert.NewCertbot(redirectService, userConfig, systemConfig, certLogger), nil
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(systemConfig *config.SystemConfig, provider *date.RealProvider) (*cert.Fake, error) {
		var certLogger *zap.Logger
		err := c.NamedResolve(&certLogger, CertificateLogger)
		if err != nil {
			return nil, err
		}
		return cert.NewFake(systemConfig, provider, cert.SubjectOrganization, cert.DefaultDuration, certLogger), nil
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *snap.SnapdHttpClient { return snap.NewSnapdHttpClient(logger) })
	if err != nil {
		return nil, err
	}

	err = c.Singleton(func(snapClient *snap.SnapdHttpClient) *snap.ChangesClient {
		return snap.NewChangesClient(snapClient, logger)
	})
	if err != nil {
		return nil, err
	}

	err = c.Singleton(func(snapClient *snap.SnapdHttpClient, deviceInfo *info.Device, systemConfig *config.SystemConfig, client *retryablehttp.Client) *snap.Server {
		return snap.NewServer(snapClient, deviceInfo, systemConfig, client, logger)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(snapServer *snap.Server, snapCli *snap.Cli, logger *zap.Logger) *event.Trigger {
		return event.New(snapServer, snapCli, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(systemConfig *config.SystemConfig, userConfig *config.UserConfig, provider *date.RealProvider, certbot *cert.Certbot, fakeCert *cert.Fake, nginxService *nginx.Nginx, eventTrigger *event.Trigger) (*cert.CertificateGenerator, error) {
		var certLogger *zap.Logger
		err := c.NamedResolve(&certLogger, CertificateLogger)
		if err != nil {
			return nil, err
		}
		return cert.New(systemConfig, userConfig, provider, certbot, fakeCert, nginxService, eventTrigger, certLogger), nil
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(certGenerator *cert.CertificateGenerator) *cron.CertificateJob {
		return cron.NewCertificateJob(certGenerator)
	})
	if err != nil {
		return nil, err
	}

	err = c.Singleton(func(userConfig *config.UserConfig, redirectService *redirect.Service, eventTrigger *event.Trigger, client *retryablehttp.Client, netInfo *network.TcpInterfaces, logger *zap.Logger) *access.PortProbe {
		return access.NewProbe(userConfig, client, logger)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(probe *access.PortProbe, userConfig *config.UserConfig, redirectService *redirect.Service, eventTrigger *event.Trigger, netInfo *network.TcpInterfaces, logger *zap.Logger) *access.ExternalAddress {
		return access.New(probe, userConfig, redirectService, eventTrigger, netInfo, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *job.SingleJobMaster { return job.NewMaster() })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(master *job.SingleJobMaster, logger *zap.Logger) *job.Worker {
		return job.NewWorker(master, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(executor *cli.ShellExecutor) *du.ShellDiskUsage {
		return du.New(executor)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(executor *cli.ShellExecutor, diskusage *du.ShellDiskUsage, snapCli *snap.Cli, snapServer *snap.Server, logger *zap.Logger, userConfig *config.UserConfig, dateProvider *date.RealProvider) *backup.Backup {
		return backup.New(backupDir, varDir, executor, diskusage, snapCli, snapServer, userConfig, dateProvider, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(job *access.ExternalAddress) *cron.ExternalAddressJob {
		return cron.NewExternalAddressJob(job)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *cron.SimpleScheduler { return &cron.SimpleScheduler{} })

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(snapd *snap.Server, userConfig *config.UserConfig, provider *date.RealProvider, backup *backup.Backup, scheduler *cron.SimpleScheduler, logger *zap.Logger) *cron.BackupJob {
		return cron.NewBackupJob(snapd, userConfig, backup, provider, scheduler, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(executor *cli.ShellExecutor, dateProvider *date.RealProvider) *cron.TimeSyncJob {
		return cron.NewTimeSyncJob(executor, dateProvider, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(job1 *cron.CertificateJob, job2 *cron.ExternalAddressJob, job3 *cron.BackupJob, job4 *cron.TimeSyncJob, userConfig *config.UserConfig) *cron.Cron {
		return cron.New([]cron.Job{job1, job2, job3, job4}, time.Minute*5, userConfig)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *installer.Installer { return installer.New() })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(systemConfig *config.SystemConfig, executor *cli.ShellExecutor, logger *zap.Logger) *storage.Storage {
		return storage.New(systemConfig, executor, 1000, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(snapService *snap.Cli, systemConfig *config.SystemConfig, executor *cli.ShellExecutor, passwordChanger *auth.SystemPasswordChanger) *auth.Service {
		return auth.New(snapService, systemConfig.DataDir(), systemConfig.AppDir(), systemConfig.ConfigDir(), executor, passwordChanger, logger)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(userConfig *config.UserConfig) *session.Cookies {
		return session.New(userConfig, logger)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(ldapService *auth.Service, nginxService *nginx.Nginx, userConfig *config.UserConfig,
		eventTrigger *event.Trigger, cookies *session.Cookies) *activation.Device {
		return activation.NewDevice(userConfig, ldapService, nginxService, eventTrigger, cookies)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() connection.InternetChecker { return connection.NewInternetChecker() })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(internetChecker connection.InternetChecker, userConfig *config.UserConfig,
		redirectService *redirect.Service, device *activation.Device, certGenerator *cert.CertificateGenerator,
		logger *zap.Logger,
	) *activation.Managed {
		return activation.NewManaged(internetChecker, userConfig, redirectService, device, certGenerator, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(internetChecker connection.InternetChecker, userConfig *config.UserConfig, device *activation.Device,
		certGenerator *cert.CertificateGenerator, logger *zap.Logger) *activation.Custom {
		return activation.NewCustom(internetChecker, userConfig, device, certGenerator, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(activationManaged *activation.Managed, activationCustom *activation.Custom) *rest.Activate {
		return rest.NewActivateBackend(activationManaged, activationCustom)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(executor *cli.ShellExecutor) *systemd.Journal { return systemd.NewJournal(executor) })

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(certGenerator *cert.CertificateGenerator, journalCtl *systemd.Journal) *rest.Certificate {
		return rest.NewCertificate(certGenerator, journalCtl)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(config *config.SystemConfig) *storage.PathChecker { return storage.NewPathChecker(config, logger) })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(systemConfig *config.SystemConfig, executor *cli.ShellExecutor, checker *storage.PathChecker) *storage.Lsblk {
		return storage.NewLsblk(systemConfig, checker, executor, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(executor *cli.ShellExecutor) *storage.FreeSpaceChecker {
		return storage.NewFreeSpaceChecker(executor)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *storage.Linker { return storage.NewLinker(logger) })
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(systemConfig *config.SystemConfig, executor *cli.ShellExecutor) *btrfs.Stats {
		return btrfs.NewStats(systemConfig, executor)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(systemConfig *config.SystemConfig, executor *cli.ShellExecutor, stats *btrfs.Stats, systemd *systemd.Control) *btrfs.Disks {
		return btrfs.NewDisks(systemConfig, executor, systemd, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(systemConfig *config.SystemConfig, freeSpaceChecker *storage.FreeSpaceChecker,
		systemd *systemd.Control, eventTrigger *event.Trigger, lsblk *storage.Lsblk,
		executor *cli.ShellExecutor, linker *storage.Linker, btrfs *btrfs.Disks, stats *btrfs.Stats) *storage.Disks {
		return storage.NewDisks(systemConfig, eventTrigger, lsblk, systemd, freeSpaceChecker, linker, executor, btrfs, stats, logger)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func() *support.LogAggregator {
		return support.NewAggregator(logger)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(aggregator *support.LogAggregator, redirectService *redirect.Service) *support.Sender {
		return support.NewSender(aggregator, redirectService)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(userConfig *config.UserConfig) *rest.Proxy {
		return rest.NewProxy(userConfig)
	})

	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(cookies *session.Cookies, userConfig *config.UserConfig) *rest.Middleware {
		return rest.NewMiddleware(cookies, userConfig, logger)
	})
	if err != nil {
		return nil, err
	}

	err = c.Singleton(func(checker *storage.PathChecker, linker *storage.Linker, systemConfig *config.SystemConfig, certGenerator *cert.CertificateGenerator, ldapService *auth.Service, nginxService *nginx.Nginx) *hook.Install {
		return hook.NewInstall(checker, linker, systemConfig, certGenerator, ldapService, nginxService, logger)
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}
