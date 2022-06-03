package storage

import (
	"github.com/syncloud/platform/config"
	"go.uber.org/zap"
	"path/filepath"
)

type PathChecker struct {
	config *config.SystemConfig
	logger *zap.Logger
}

type Checker interface {
	ExternalDiskLinkExists() bool
}

func NewPathChecker(config *config.SystemConfig, logger *zap.Logger) *PathChecker {
	return &PathChecker{
		config: config,
		logger: logger,
	}
}

func (c *PathChecker) ExternalDiskLinkExists() bool {
	realLinkPath, err := filepath.EvalSymlinks(c.config.DiskLink())
	if err != nil {
		c.logger.Error("cannot read disk link", zap.String("name", c.config.DiskLink()))
		return false
	}
	c.logger.Info("real link", zap.String("path", realLinkPath))

	externalDiskPath := c.config.ExternalDiskDir()
	c.logger.Info("external disk", zap.String("path", externalDiskPath))

	linkExists := realLinkPath == externalDiskPath
	c.logger.Info("link", zap.Bool("exists", linkExists))

	return linkExists
}
