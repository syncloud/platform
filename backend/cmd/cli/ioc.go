package main

import (
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/ioc"
)

func Init(userConfig string, systemConfig string) {
	ioc.Init(userConfig, systemConfig, backup.Dir, backup.VarDir)
}
