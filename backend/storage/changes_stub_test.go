package storage

import snapmodel "github.com/syncloud/platform/snap/model"

type ChangesStub struct {
	running bool
}

func (c *ChangesStub) Changes() (*snapmodel.InstallerStatus, error) {
	return &snapmodel.InstallerStatus{IsRunning: c.running}, nil
}
