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
	last *time.Time
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

func (c *UserConfigStub) GetBackupAppTime(_ string, _ string) *time.Time {
	return c.last
}

func (c *UserConfigStub) SetBackupAppTime(_ string, _ string, last time.Time) {
	c.last = &last
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

func TestRun_Disabled(t *testing.T) {
	snapd := &SnapdStub{}
	config := &UserConfigStub{auto: "no"}
	backup := &BackupStub{}
	timeProvider := &ProviderStub{}
	job := NewBackupJob(snapd, config, backup, timeProvider, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backup.created)
	assert.False(t, backup.restored)
}

func TestRun_Daily_LessThanOneDay_SinceLastRun(t *testing.T) {
	snapd := &SnapdStub{}
	now := time.Date(2009, 1, 1, 1, 0, 0, 0, time.UTC)
	timeProvider := &ProviderStub{now}
	last := now.Add(-23 * time.Hour)
	config := &UserConfigStub{auto: "backup", day: 0, hour: 1, last: &last}
	backup := &BackupStub{}
	job := NewBackupJob(snapd, config, backup, timeProvider, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backup.created)
	assert.False(t, backup.restored)
}

func TestRun_Daily_MoreThanOneDay_SinceLastRun(t *testing.T) {
	snapd := &SnapdStub{}
	now := time.Date(2022, 11, 21, 1, 0, 0, 0, time.UTC)
	timeProvider := &ProviderStub{now}
	last := now.Add(-25 * time.Hour)
	config := &UserConfigStub{auto: "backup", day: 0, hour: 1, last: &last}
	backup := &BackupStub{}
	job := NewBackupJob(snapd, config, backup, timeProvider, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.True(t, backup.created)
	assert.False(t, backup.restored)
}

func TestRun_Monday_LessThanOneWeek_SinceLastRun(t *testing.T) {
	snapd := &SnapdStub{}
	monday := time.Date(2022, 11, 21, 1, 0, 0, 0, time.UTC)
	timeProvider := &ProviderStub{monday.Add(24 * time.Hour)}
	config := &UserConfigStub{auto: "backup", day: 1, hour: 1, last: &monday}
	backup := &BackupStub{}
	job := NewBackupJob(snapd, config, backup, timeProvider, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backup.created)
	assert.False(t, backup.restored)
}

func TestRun_Monday_MoreThanOneWeek_SinceLastRun(t *testing.T) {
	snapd := &SnapdStub{}
	monday := time.Date(2022, 11, 21, 1, 0, 0, 0, time.UTC)
	timeProvider := &ProviderStub{monday.AddDate(0, 0, 8)}
	config := &UserConfigStub{auto: "backup", day: 1, hour: 1, last: &monday}
	backup := &BackupStub{}
	job := NewBackupJob(snapd, config, backup, timeProvider, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backup.created)
	assert.False(t, backup.restored)
}

func TestRun_Monday_FirstTime_NonMonday_NotRun(t *testing.T) {
	snapd := &SnapdStub{}
	monday := time.Date(2022, 11, 21, 1, 0, 0, 0, time.UTC)
	timeProvider := &ProviderStub{monday.AddDate(0, 0, 1)}
	config := &UserConfigStub{auto: "backup", day: 1, hour: 1, last: nil}
	backup := &BackupStub{}
	job := NewBackupJob(snapd, config, backup, timeProvider, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backup.created)
	assert.False(t, backup.restored)
}

func TestRun_Monday_FirstTime_Monday_Run(t *testing.T) {
	snapd := &SnapdStub{}
	monday := time.Date(2022, 11, 21, 1, 0, 0, 0, time.UTC)
	timeProvider := &ProviderStub{monday.AddDate(0, 0, 1)}
	config := &UserConfigStub{auto: "backup", day: 1, hour: 1, last: nil}
	backup := &BackupStub{}
	job := NewBackupJob(snapd, config, backup, timeProvider, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.True(t, backup.created)
	assert.False(t, backup.restored)
}
