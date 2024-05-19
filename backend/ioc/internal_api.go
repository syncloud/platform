package ioc

import (
	"github.com/golobby/container/v3"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/storage"
	"github.com/syncloud/platform/systemd"
)

func InitInternalApi(userConfig string, systemConfig string, backupDir string, varDir string, network string, address string) (container.Container, error) {
	c, err := Init(userConfig, systemConfig, backupDir, varDir)
	if err != nil {
		return nil, err
	}
	err = c.Singleton(func(
		userConfig *config.UserConfig,
		storage *storage.Storage,
		systemd *systemd.Control,
		middleware *rest.Middleware,
		authelia *auth.Authelia,
	) *rest.Api {
		return rest.NewApi(userConfig, storage, systemd, middleware, network, address, authelia, logger)
	})
	if err != nil {
		return nil, err
	}

	err = c.Singleton(func(
		api *rest.Api,
	) []Service {
		return []Service{
			api,
		}
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
