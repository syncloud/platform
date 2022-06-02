package storage

import (
	"github.com/syncloud/platform/config"
	"go.uber.org/zap"
)

type PathChecker struct {
	config *config.SystemConfig
	logger *zap.Logger
}

func NewPathChecker(config *config.SystemConfig, logger *zap.Logger) *PathChecker {
	return &PathChecker{
		config: config,
		logger: logger,
	}
}

func (c *PathChecker) ExternalDiskLinkExists() bool {
	real_link_path = path.realpath(self.platform_config.get_disk_link())
	self.log.info('real link path: {0}'.format(real_link_path))

	external_disk_path = self.platform_config.get_external_disk_dir()
	self.log.info('external disk path: {0}'.format(external_disk_path))

	link_exists = real_link_path == external_disk_path
	self.log.info('link exists: {0}'.format(link_exists))

	return link_exists
}
