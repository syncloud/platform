package ioc

import (
	"github.com/golobby/container/v3"
	"github.com/syncloud/platform/access"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/cron"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/hardware/lcd"
	"github.com/syncloud/platform/health"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/network"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/session"
	"github.com/syncloud/platform/snap"
	"github.com/syncloud/platform/storage"
	"github.com/syncloud/platform/support"
	"github.com/syncloud/platform/systemd"
	"github.com/syncloud/platform/timezone"
)

func InitPublicApi(userConfig string, systemConfig string, backupDir string, varDir string, net string, address string) (container.Container, error) {
	c, err := Init(userConfig, systemConfig, backupDir, varDir)
	if err != nil {
		return nil, err
	}

	err = c.Singleton(func(master *job.SingleJobMaster, backupService *backup.Backup, eventTrigger *event.Trigger, worker *job.Worker,
		redirectService *redirect.Service, snapdUpgrader *snap.Snapd, storageService *storage.Storage,
		id *identification.Parser, activate *rest.Activate, userConfig *config.UserConfig, redirectConfig *config.Redirect, cert *rest.Certificate,
		externalAddress *access.ExternalAddress, snapd *snap.Server, disks *storage.Disks, journalCtl *systemd.Journal,
		executor *cli.ShellExecutor, iface *network.TcpInterfaces, sender *support.Sender,
		proxy *rest.Proxy, customProxy *rest.CustomProxy, middleware *rest.Middleware, userManager *auth.UserManager, groupManager *auth.GroupManager, cookies *session.Cookies,
		changesClient *snap.ChangesClient,
		oidcService *auth.OIDCService, authelia *auth.Authelia, totp *auth.TOTP,
		tz *timezone.Applier,
		healthService *health.Health,
	) *rest.Backend {
		return rest.NewBackend(master, backupService, eventTrigger, worker, redirectService,
			snapdUpgrader, storageService, id, activate, userConfig, redirectConfig, cert, externalAddress,
			snapd, disks, journalCtl, executor, iface, sender, proxy, customProxy,
			userManager, groupManager, middleware, cookies, net, address, changesClient,
			oidcService, authelia, totp, tz, healthService, logger)
	})
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(
		cronService *cron.Cron,
		backupService *backup.Backup,
		cookies *session.Cookies,
		backend *rest.Backend,
		lcdDisplay *lcd.Display,
	) []Service {
		return []Service{
			cronService,
			backupService,
			cookies,
			backend,
			lcdDisplay,
		}
	})
	if err != nil {
		return nil, err
	}
	return c, err
}
