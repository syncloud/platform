package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

func TestStartupTimeSync_RunsOnceWhenSynced(t *testing.T) {
	job := &JobStub{}
	assert.NoError(t, NewStartupTimeSync(job, time.Second, time.Millisecond, log.Default()).Start())
	assert.Equal(t, 1, job.runs)
}

func TestStartupTimeSync_RetriesUntilNetworkIsUp(t *testing.T) {
	job := &JobStub{failures: 3}
	assert.NoError(t, NewStartupTimeSync(job, time.Second, time.Millisecond, log.Default()).Start())
	assert.Equal(t, 4, job.runs)
}

func TestStartupTimeSync_GivesUpAfterTimeout(t *testing.T) {
	job := &JobStub{failures: 1000}
	start := time.Now()
	assert.NoError(t, NewStartupTimeSync(job, 50*time.Millisecond, 10*time.Millisecond, log.Default()).Start())
	assert.Less(t, time.Since(start), 100*time.Millisecond)
	assert.GreaterOrEqual(t, job.runs, 4)
	assert.Less(t, job.runs, 1000)
}
