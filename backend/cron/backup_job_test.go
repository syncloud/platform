package cron

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/snap/model"
	"testing"
	"time"
)

type SnapdStub struct {
	app string
}

func (s *SnapdStub) InstalledUserApps() ([]model.SyncloudApp, error) {
	return []model.SyncloudApp{{Name: s.app}}, nil
}

type UserConfigStub struct {
	auto     string
	day      int
	hour     int
	lastTime time.Time
	lastMode string
	lastApp  string
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
	return c.lastTime
}

func (c *UserConfigStub) SetBackupAppTime(app string, mode string, last time.Time) {
	c.lastTime = last
	c.lastMode = mode
	c.lastApp = app

}

type BackupStub struct {
	err      error
	created  bool
	restored bool
	list     []backup.File
}

func (b *BackupStub) Create(_ string) error {
	if b.err != nil {
		return b.err
	}
	b.created = true
	return nil
}

func (b *BackupStub) Restore(_ string) error {
	if b.err != nil {
		return b.err
	}
	b.restored = true
	return nil
}

func (b *BackupStub) List() ([]backup.File, error) {
	return b.list, nil
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
	backuper := &BackupStub{}
	timeProvider := &ProviderStub{}
	scheduler := &SchedulerStub{}
	job := NewBackupJob(snapd, config, backuper, timeProvider, scheduler, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backuper.created)
	assert.False(t, backuper.restored)
}

func TestRun_Backup(t *testing.T) {
	snapd := &SnapdStub{app: "app1"}
	config := &UserConfigStub{auto: "backup"}
	backuper := &BackupStub{}
	timeProvider := &ProviderStub{}
	scheduler := &SchedulerStub{run: true}
	job := NewBackupJob(snapd, config, backuper, timeProvider, scheduler, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.True(t, backuper.created)
	assert.False(t, backuper.restored)
	assert.Equal(t, "backup", config.lastMode)
	assert.Equal(t, "app1", config.lastApp)
}

func TestRun_Backup_Failed(t *testing.T) {
	snapd := &SnapdStub{app: "app1"}
	config := &UserConfigStub{auto: "backup"}
	backuper := &BackupStub{err: fmt.Errorf("expected error")}
	timeProvider := &ProviderStub{}
	scheduler := &SchedulerStub{run: true}
	job := NewBackupJob(snapd, config, backuper, timeProvider, scheduler, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backuper.created)
	assert.False(t, backuper.restored)
	assert.Equal(t, "", config.lastMode)
	assert.Equal(t, "", config.lastApp)
}

func TestRun_Restore(t *testing.T) {
	snapd := &SnapdStub{app: "app1"}
	config := &UserConfigStub{auto: "restore"}
	backuper := &BackupStub{
		list: []backup.File{
			{Path: "/data/platform/backup", File: "app1-2020-0514-061314.tar.gz", App: "app1"},
		},
	}
	timeProvider := &ProviderStub{}
	scheduler := &SchedulerStub{run: true}
	job := NewBackupJob(snapd, config, backuper, timeProvider, scheduler, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backuper.created)
	assert.True(t, backuper.restored)
	assert.Equal(t, "restore", config.lastMode)
	assert.Equal(t, "app1", config.lastApp)
}

func TestRun_Restore_Failed(t *testing.T) {
	snapd := &SnapdStub{app: "app1"}
	config := &UserConfigStub{auto: "restore"}
	backuper := &BackupStub{err: fmt.Errorf("expected error")}
	timeProvider := &ProviderStub{}
	scheduler := &SchedulerStub{run: true}
	job := NewBackupJob(snapd, config, backuper, timeProvider, scheduler, log.Default())
	err := job.Run()

	assert.Nil(t, err)
	assert.False(t, backuper.created)
	assert.False(t, backuper.restored)
	assert.Equal(t, "", config.lastMode)
	assert.Equal(t, "", config.lastApp)
}

func TestBackupJob_LatestBackup(t *testing.T) {
	snapd := &SnapdStub{app: "app1"}
	config := &UserConfigStub{auto: "restore"}
	backuper := &BackupStub{
		list: []backup.File{
			{Path: "/data/platform/backup", File: "files-2020-0514-061314.tar.gz", App: "files"},
			{Path: "/data/platform/backup", File: "mail-2020-0910-110840.tar.gz", App: "mail"},
			{Path: "/data/platform/backup", File: "rocketchat-2019-0710-204410.tar.gz", App: "rocketchat"},
			{Path: "/data/platform/backup", File: "rocketchat-2020-0107-230800.tar.gz", App: "rocketchat"},
			{Path: "/data/platform/backup", File: "rocketchat-2020-0910-111456.tar.gz", App: "rocketchat"},
			{Path: "/data/platform/backup", File: "rocketchat-2022-0129-231530.tar.gz", App: "rocketchat"},
			{Path: "/data/platform/backup", File: "rocketchat-2022-0128-231530.tar.gz", App: "rocketchat"},
			{Path: "/data/platform/backup", File: "rocketchat-2022-0129-231529.tar.gz", App: "rocketchat"},
			{Path: "/data/platform/backup", File: "rocketchat-2020-1230-100714.tar.gz", App: "rocketchat"},
			{Path: "/data/platform/backup", File: "mail-2021-0109-130507.tar.gz", App: "mail"},
		},
	}
	timeProvider := &ProviderStub{}
	scheduler := &SchedulerStub{run: true}
	job := NewBackupJob(snapd, config, backuper, timeProvider, scheduler, log.Default())
	latest, err := job.LatestBackup("rocketchat")
	assert.Nil(t, err)
	assert.Equal(t, "rocketchat-2022-0129-231530.tar.gz", latest)
}
