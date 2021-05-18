package event

import (
	"github.com/syncloud/platform/snap"
	"log"
	"os/exec"
)

type Trigger struct {
	snap *snap.Snapd
}

func New() *Trigger {
	return &Trigger{
		snap: snap.New(),
	}
}

func (storage *Trigger) RunEventOnAllApps(event string) error {

	snaps, err := storage.snap.ListAllApps()
	if err != nil {
		log.Printf("snap info failed: %v", err)
		return err
	}
	for _, snap := range snaps {
		found, app := snap.FindApp(event)
		if found {
			var cmd = app.RunCommand()
			log.Println("Running: ", cmd)
			_, err := exec.Command("snap", "run", cmd).CombinedOutput()
			if err != nil {
				log.Printf("snap run failed: %v", err)
				return err
			}
		}
	}
	return nil
}
