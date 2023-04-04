package cron

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
	"time"
)

type ExecutorStub struct {
	called int
}

func (e *ExecutorStub) CombinedOutput(_ string, _ ...string) ([]byte, error) {
	e.called++
	return []byte("test output"), nil
}

type DateProviderStub struct {
	now time.Time
}

func (d *DateProviderStub) Now() time.Time {
	return d.now
}

func TestTimeSyncJob_RunOnce_24h(t *testing.T) {
	executor := &ExecutorStub{}
	date := &DateProviderStub{}
	job := NewTimeSyncJob(executor, date, log.Default())
	date.now = time.Now()
	err := job.Run()
	assert.NoError(t, err)
	date.now = date.now.Add(time.Hour * 24)
	err = job.Run()
	err = job.Run()
	assert.NoError(t, err)
	assert.Equal(t, 3, executor.called)
}

func TestTimeSyncJob_RunTwice_48h(t *testing.T) {
	executor := &ExecutorStub{}
	date := &DateProviderStub{}
	job := NewTimeSyncJob(executor, date, log.Default())
	date.now = time.Now()
	err := job.Run()
	assert.NoError(t, err)
	date.now = date.now.Add(time.Hour * 48)
	err = job.Run()
	err = job.Run()
	assert.NoError(t, err)
	assert.Equal(t, 6, executor.called)
}
