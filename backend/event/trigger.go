package event

import (
	"github.com/syncloud/platform/snap"
	"log"
	"os/exec"
)

type Trigger struct {
	snapd *snap.Snapd
}

func New(snapd *snap.Snapd) *Trigger {
	return &Trigger{
		snapd: snapd,
	}
}

func (t *Trigger) RunAccessChangeEvent() error {
	return t.RunEventOnAllApps("access-change")
}

func (t *Trigger) RunEventOnAllApps(event string) error {

	snaps, err := t.snapd.InstalledSnaps()
	if err != nil {
		log.Printf("snap info failed: %v", err)
		return err
	}
	for _, info := range snaps {
		found, app := info.FindApp(event)
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
