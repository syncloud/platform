package event

import (
	"github.com/syncloud/platform/snapd"
	"log"
	"os/exec"
)

type Trigger struct {
	snap *snapd.Snapd
}

func New() *Trigger {
	return &Trigger{
		snap: snapd.New(),
	}
}

func (storage *Trigger) RunEventOnAllAps(event string) error {

	snaps, err := storage.snap.ListAllApps()
	if err != nil {
		log.Printf("snapd info failed: %v", err)
		return err
	}
	for _, snap := range snaps {
		found, app := snap.FindApp(event)
		if found {
			var cmd = app.RunCommand()
			log.Println("Running: ", cmd)
			_, err := exec.Command("snap", "run", cmd).CombinedOutput()
			if err != nil {
				log.Printf("snapd run failed: %v", err)
				return err
			}
		}
	}
	return nil
}
