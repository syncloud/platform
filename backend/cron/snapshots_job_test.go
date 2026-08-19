package cron

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

type SnapshotsCleanerStub struct {
	disabled  int
	forgotten int
	err       error
}

func (c *SnapshotsCleanerStub) DisableAutomatic() error {
	c.disabled++
	return c.err
}

func (c *SnapshotsCleanerStub) ForgetAuto() error {
	c.forgotten++
	return nil
}

func TestSnapshotsJob_DisablesAndForgets(t *testing.T) {
	cleaner := &SnapshotsCleanerStub{}
	job := NewSnapshotsJob(cleaner, log.Default())

	assert.NoError(t, job.Run())
	assert.Equal(t, 1, cleaner.disabled)
	assert.Equal(t, 1, cleaner.forgotten)
}

func TestSnapshotsJob_DoesNotForgetWhenDisableFails(t *testing.T) {
	cleaner := &SnapshotsCleanerStub{err: fmt.Errorf("snapd is down")}
	job := NewSnapshotsJob(cleaner, log.Default())

	assert.Error(t, job.Run())
	assert.Equal(t, 0, cleaner.forgotten)
}
