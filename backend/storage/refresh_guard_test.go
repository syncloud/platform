package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/job"
	"testing"
)

func TestRefreshGuard_Idle(t *testing.T) {
	guard := NewRefreshGuard(&ChangesStub{}, &JobsStub{status: job.NewStatus("", job.Idle)})
	assert.Nil(t, guard.Check())
}

func TestRefreshGuard_AppRefresh(t *testing.T) {
	guard := NewRefreshGuard(&ChangesStub{running: true}, &JobsStub{status: job.NewStatus("", job.Idle)})
	assert.EqualError(t, guard.Check(), "an app install or update is in progress, try again when it finishes")
}

func TestRefreshGuard_PlatformJob(t *testing.T) {
	guard := NewRefreshGuard(&ChangesStub{}, &JobsStub{status: job.NewStatus("storage.activate.disks", job.Busy)})
	assert.EqualError(t, guard.Check(), "storage.activate.disks is in progress, try again when it finishes")
}
