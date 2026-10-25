package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

var syncNow = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func newTimeSyncJob(offsets map[string]time.Duration) (*TimeSyncJob, *TimeSourceStub, *ClockStub) {
	source := &TimeSourceStub{offsets: offsets}
	clock := &ClockStub{}
	job := NewTimeSyncJob(source, clock, []string{"a", "b", "c"}, &DateProviderStub{now: syncNow}, log.Default())
	return job, source, clock
}

func TestTimeSyncJob_LeavesSmallDriftToNtpd(t *testing.T) {
	job, source, clock := newTimeSyncJob(map[string]time.Duration{"a": 30 * time.Second})
	assert.NoError(t, job.Run())
	assert.Equal(t, []string{"a"}, source.asked)
	assert.Empty(t, clock.set)
}

func TestTimeSyncJob_SetsClockWhenFarOff(t *testing.T) {
	offset := 56 * 365 * 24 * time.Hour
	job, _, clock := newTimeSyncJob(map[string]time.Duration{"a": offset})
	assert.NoError(t, job.Run())
	assert.Equal(t, []time.Time{syncNow.Add(offset)}, clock.set)
}

func TestTimeSyncJob_SetsClockWhenAhead(t *testing.T) {
	job, _, clock := newTimeSyncJob(map[string]time.Duration{"a": -2 * time.Hour})
	assert.NoError(t, job.Run())
	assert.Equal(t, []time.Time{syncNow.Add(-2 * time.Hour)}, clock.set)
}

func TestTimeSyncJob_TriesNextServer(t *testing.T) {
	job, source, clock := newTimeSyncJob(map[string]time.Duration{"c": time.Hour})
	assert.NoError(t, job.Run())
	assert.Equal(t, []string{"a", "b", "c"}, source.asked)
	assert.Equal(t, []time.Time{syncNow.Add(time.Hour)}, clock.set)
}

func TestTimeSyncJob_FailsWhenNoServerAnswers(t *testing.T) {
	job, _, clock := newTimeSyncJob(map[string]time.Duration{})
	err := job.Run()
	assert.ErrorContains(t, err, "a: i/o timeout")
	assert.ErrorContains(t, err, "c: i/o timeout")
	assert.Empty(t, clock.set)
}
