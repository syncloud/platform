package main

import (
	"github.com/golobby/container/v3"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/ioc"
)

func Init(userConfig string, systemConfig string) (container.Container, error) {
	return ioc.Init(userConfig, systemConfig, backup.Dir, backup.VarDir)
}
