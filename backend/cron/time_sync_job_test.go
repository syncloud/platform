package cron

import (
	"fmt"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
	"time"
)

type ExecutorStub struct {
	called   int
	commands []string
	failOn   string
}

func (e *ExecutorStub) CombinedOutput(name string, arg ...string) ([]byte, error) {
	e.called++
	command := strings.Join(append([]string{name}, arg...), " ")
	e.commands = append(e.commands, command)
	if e.failOn != "" && strings.Contains(command, e.failOn) {
		return []byte("failed"), fmt.Errorf("exit status 124")
	}
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

func TestTimeSyncJob_BoundsNtpd(t *testing.T) {
	executor := &ExecutorStub{}
	job := NewTimeSyncJob(executor, &DateProviderStub{now: time.Now()}, log.Default())
	assert.NoError(t, job.Run())
	assert.Equal(t, []string{"service ntp stop", "timeout 60 ntpd -gq", "service ntp start"}, executor.commands)
}

func TestTimeSyncJob_StartsNtpAndRetries_WhenSyncFails(t *testing.T) {
	executor := &ExecutorStub{failOn: "ntpd"}
	date := &DateProviderStub{now: time.Now()}
	job := NewTimeSyncJob(executor, date, log.Default())
	assert.Error(t, job.Run())
	assert.Equal(t, "service ntp start", executor.commands[2])
	executor.failOn = ""
	assert.NoError(t, job.Run())
	assert.Equal(t, 6, executor.called)
}
