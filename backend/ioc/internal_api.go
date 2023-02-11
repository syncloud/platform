package ioc

import (
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/info"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/storage"
	"github.com/syncloud/platform/systemd"
)

func InitInternalApi(userConfig string, systemConfig string, backupDir string, varDir string, network string, address string) {
	Init(userConfig, systemConfig, backupDir, varDir)
	Singleton(func(device *info.Device, userConfig *config.UserConfig, storage *storage.Storage,
		systemd *systemd.Control, middleware *rest.Middleware) *rest.Api {
		return rest.NewApi(device, userConfig, storage, systemd, middleware, network, address, logger)
	})

	Singleton(func(
		api *rest.Api,
	) []Service {
		return []Service{
			api,
		}
	})
}
