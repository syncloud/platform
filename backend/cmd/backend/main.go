package main

import (
	"fmt"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/redirect"
	"log"
	"os"

	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/storage"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: ", os.Args[0], "/path.sock")
		return
	}

	os.Remove(os.Args[1])
	Backend().Start(os.Args[1])

}

func Backend() *rest.Backend {
	master := job.NewMaster()
	backupService := backup.NewDefault()
	eventTrigger := event.New()
	installerService := installer.New()
	storageService := storage.New()
	oldConfig := fmt.Sprintf("%s/user_platform.cfg", os.Getenv("SNAP_COMMON"))
	conf := fmt.Sprintf("%s/platform.db", os.Getenv("SNAP_COMMON"))
	configuration := config.New(conf, oldConfig)
	redirectService := redirect.New(configuration)
	worker := job.NewWorker(master, backupService, installerService, storageService)

	return rest.NewBackend(master, backupService, eventTrigger, worker, redirectService)

}
