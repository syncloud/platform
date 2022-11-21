package cron

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/snap/model"
	"testing"
)

type SnapdStub struct {
}

func (s *SnapdStub) InstalledUserApps() ([]model.SyncloudApp, error) {
    return []model.SyncloudApp{{ Name: "app1" }}, nil
}

type UserConfigStub struct {
 auto string
 day int
 hour int
 last *int64
}

func (c *UserConfigStub) GetBackupAuto() string {
 return c.auto
}

func (c *UserConfigStub)	GetBackupAutoDay() int {
 return c.day
}

func (c *UserConfigStub)	GetBackupAutoHour() int {
 return c.hour
}

func (c *UserConfigStub) GetBackupAppTime(app string, mode string) *int64 {
 return c.last
}

func (c *UserConfigStub)	SetBackupAppTime(app string, mode string, last int64) {
 c.last = &last
}


func TestRun(t *testing.T) {
	snapd := &SnapdStub{}
	config := &UserConfigStub{}
	job := NewBackupJob(snapd, config, log.Default())
	
	err := job.Run()
	assert.Nil(t, err)
}
