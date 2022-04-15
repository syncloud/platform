package cron

import (
	"github.com/syncloud/platform/access"
)

type ExternalAddressJob struct {
	externalAddress *access.ExternalAddress
}

func NewExternalAddressJob(externalAddress *access.ExternalAddress) *ExternalAddressJob {
	return &ExternalAddressJob{
		externalAddress: externalAddress,
	}
}

func (j *ExternalAddressJob) Run() error {
	j.externalAddress.Sync()
	return nil
}
