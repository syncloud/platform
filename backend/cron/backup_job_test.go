package cron

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/snap/model"
	"testing"
	"time"
)

type SnapdStub struct {
}

func (s *SnapdStub) InstalledUserApps() ([]model.SyncloudApp, error) {
	return []model.SyncloudApp{{Name: "app1"}}, nil
}

type UserConfigStub struct {
	auto string
	day  int
	hour int
	last time.Time
}

func (c *UserConfigStub) GetBackupAuto() string {
	return c.auto
}

func (c *UserConfigStub) GetBackupAutoDay() int {
	return c.day
}

func (c *UserConfigStub) GetBackupAutoHour() int {
	return c.hour
}

func (c *UserConfigStub) GetBackupAppTime(_ string, _ string) time.Time {
	return c.last
}

func (c *UserConfigStub) SetBackupAppTime(_ string, _ string, last time.Time) {
	c.last = last
}

type BackupStub struct {
	created  bool
	restored bool
}

func (b *BackupStub) Create(_ string) error {
	b.created = true
	return nil
}

func (b *BackupStub) Restore(_ string) error {
	b.restored = true
	return nil
}

type ProviderStub struct {
	now time.Time
}

func (p ProviderStub) Now() time.Time {
	return p.now
}

type SchedulerStub struct {
	run bool
}

func (s *SchedulerStub) ShouldRun(_ int, _ int, _ time.Time, _ time.Time) bool {
	return s.run
}

func TestRun_Disabled(t *testing.T) {
	snapd := &SnapdStub{}
	config := &UserConfigStub{auto: "no"}
	backup := &BackupStub{}
	timeProvider := &ProviderStub{}
	scheduler := &SchedulerStub{}
	job := NewBackupJob(snapd, config, backup, timeProvider, scheduler, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backup.created)
	assert.False(t, backup.restored)
}
