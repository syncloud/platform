package event

import (
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
)

type Trigger struct {
	snapServer SnapServer
	snapCli    SnapRunner
	logger     *zap.Logger
}

type SnapServer interface {
	Snaps() ([]model.Snap, error)
}

type SnapRunner interface {
	Run(name string) error
}

func New(snapServer SnapServer, snapCli SnapRunner, logger *zap.Logger) *Trigger {
	return &Trigger{
		snapServer: snapServer,
		snapCli:    snapCli,
		logger:     logger,
	}
}

func (t *Trigger) RunAccessChangeEvent() error {
	return t.RunEventOnAllApps("access-change")
}

func (t *Trigger) RunDiskChangeEvent() error {
	return t.RunEventOnAllApps("storage-change")
}

func (t *Trigger) RunEventOnAllApps(command string) error {

	snaps, err := t.snapServer.Snaps()
	if err != nil {
		return err
	}
	for _, app := range snaps {
		cmd := app.FindCommand(command)
		if cmd != nil {
			err = t.snapCli.Run(cmd.FullName())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
