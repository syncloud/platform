package event

import (
	"github.com/syncloud/platform/executor"
	"github.com/syncloud/platform/snap/model"
	"log"
)

type Trigger struct {
	snapd    Snapd
	executor executor.Executor
}

type Snapd interface {
	InstalledSnaps() ([]model.Snap, error)
}

func New(snapd Snapd, executor executor.Executor) *Trigger {
	return &Trigger{
		snapd:    snapd,
		executor: executor,
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
			_, err := t.executor.CommandOutput("snap", "run", cmd)
			if err != nil {
				log.Printf("snap run failed: %v", err)
				return err
			}
		}
	}
	return nil
}
