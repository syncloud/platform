package storage

import (
	"fmt"
	"github.com/syncloud/platform/job"
	snapmodel "github.com/syncloud/platform/snap/model"
)

type GuardChanges interface {
	Changes() (*snapmodel.InstallerStatus, error)
}

type GuardJobs interface {
	Status() job.Status
}

type RefreshGuard struct {
	changes GuardChanges
	jobs    GuardJobs
}

func NewRefreshGuard(changes GuardChanges, jobs GuardJobs) *RefreshGuard {
	return &RefreshGuard{
		changes: changes,
		jobs:    jobs,
	}
}

func (g *RefreshGuard) Check() error {
	status := g.jobs.Status()
	if status.Status != "Idle" {
		return fmt.Errorf("%s is in progress, try again when it finishes", status.Name)
	}
	changes, err := g.changes.Changes()
	if err != nil {
		return err
	}
	if changes.IsRunning {
		return fmt.Errorf("an app install or update is in progress, try again when it finishes")
	}
	return nil
}
