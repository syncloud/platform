package cron

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
	"testing"
	"time"
)

func at(day int, hour int, minute int) time.Time {
	return time.Date(2026, 11, day, hour, minute, 0, 0, time.UTC)
}

func TestStorageCheckJob_OnlyAtCheckHour(t *testing.T) {
	checks := &StorageChecksStub{schedule: model.CheckSchedule{Enabled: true}, scrubAvailable: true}
	config := &StorageCheckConfigStub{times: map[string]time.Time{}}
	job := NewStorageCheckJob(checks, config, &DateProviderStub{now: at(1, 2, 55)}, log.Default())
	assert.Nil(t, job.Run())
	assert.Equal(t, 0, checks.scrubs)
	assert.Equal(t, 0, checks.smartTests)
}

func TestStorageCheckJob_Disabled(t *testing.T) {
	checks := &StorageChecksStub{scrubAvailable: true}
	config := &StorageCheckConfigStub{times: map[string]time.Time{}}
	job := NewStorageCheckJob(checks, config, &DateProviderStub{now: at(1, 3, 0)}, log.Default())
	assert.Nil(t, job.Run())
	assert.Equal(t, 0, checks.scrubs)
	assert.Equal(t, 0, checks.smartTests)
}

func TestStorageCheckJob_MonthlyScrubWeeklySmart(t *testing.T) {
	checks := &StorageChecksStub{schedule: model.CheckSchedule{Enabled: true}, scrubAvailable: true}
	config := &StorageCheckConfigStub{times: map[string]time.Time{}}
	provider := &DateProviderStub{}
	job := NewStorageCheckJob(checks, config, provider, log.Default())

	for day := 1; day <= 30; day++ {
		for minute := 0; minute < 60; minute += 5 {
			provider.now = at(day, 3, minute)
			assert.Nil(t, job.Run())
		}
	}
	assert.Equal(t, 1, checks.scrubs)
	assert.Equal(t, 5, checks.smartTests)

	provider.now = time.Date(2026, 12, 1, 3, 0, 0, 0, time.UTC)
	assert.Nil(t, job.Run())
	assert.Equal(t, 2, checks.scrubs)
}

func TestStorageCheckJob_ScrubSkippedWithoutArray(t *testing.T) {
	checks := &StorageChecksStub{schedule: model.CheckSchedule{Enabled: true}}
	config := &StorageCheckConfigStub{times: map[string]time.Time{}}
	job := NewStorageCheckJob(checks, config, &DateProviderStub{now: at(1, 3, 0)}, log.Default())
	assert.Nil(t, job.Run())
	assert.Equal(t, 0, checks.scrubs)
	assert.True(t, config.times[StorageCheckScrub].IsZero())
}

func TestStorageCheckJob_RetriesWhenBusy(t *testing.T) {
	checks := &StorageChecksStub{schedule: model.CheckSchedule{Enabled: true}, scrubAvailable: true, scrubErr: fmt.Errorf("an app install or update is in progress")}
	config := &StorageCheckConfigStub{times: map[string]time.Time{}}
	provider := &DateProviderStub{now: at(1, 3, 0)}
	job := NewStorageCheckJob(checks, config, provider, log.Default())

	assert.EqualError(t, job.Run(), "an app install or update is in progress")
	assert.True(t, config.times[StorageCheckScrub].IsZero())

	checks.scrubErr = nil
	provider.now = at(1, 3, 5)
	assert.Nil(t, job.Run())
	assert.Equal(t, 1, checks.scrubs)
	assert.Equal(t, at(1, 3, 5), config.times[StorageCheckScrub])
}
