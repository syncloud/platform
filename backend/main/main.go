package main

import (
	"log"
	"os"

	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/rest"
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
	backup := backup.NewDefault()
	worker := job.NewWorker(master, backup)

	return rest.NewBackend(master, backup, worker)

}
