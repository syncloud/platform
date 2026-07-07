package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/snap/model"
)

type InstallerInfoStub struct {
	info *model.InstallerInfo
	err  error
}

func (s *InstallerInfoStub) Installer() (*model.InstallerInfo, error) {
	return s.info, s.err
}

type UpgraderStub struct {
	called int
}

func (u *UpgraderStub) Upgrade() error {
	u.called++
	return nil
}

type JobMasterStub struct {
	offered []string
	busy    bool
}

func (m *JobMasterStub) Offer(name string, j job.Job) error {
	if m.busy {
		return assert.AnError
	}
	m.offered = append(m.offered, name)
	return j()
}

func atHour(hour int) time.Time {
	return time.Date(2026, 7, 7, hour, 0, 0, 0, time.UTC)
}

func newJob(store, installed string) (*SnapdUpgradeJob, *UpgraderStub, *JobMasterStub, *DateProviderStub) {
	info := &InstallerInfoStub{info: &model.InstallerInfo{StoreVersion: store, InstalledVersion: installed}}
	upgrader := &UpgraderStub{}
	master := &JobMasterStub{}
	provider := &DateProviderStub{}
	j := NewSnapdUpgradeJob(info, upgrader, master, &SimpleScheduler{}, provider, log.Default())
	return j, upgrader, master, provider
}

func TestSnapdUpgradeJob_UpgradesWhenVersionsDiffer(t *testing.T) {
	j, upgrader, master, provider := newJob("2.60.0", "2.59.5")
	provider.now = atHour(SnapdUpgradeHour)

	err := j.Run()

	assert.NoError(t, err)
	assert.Equal(t, 1, upgrader.called)
	assert.Equal(t, []string{"installer.upgrade"}, master.offered)
}

func TestSnapdUpgradeJob_SkipsWhenUpToDate(t *testing.T) {
	j, upgrader, master, provider := newJob("2.59.5\n", "2.59.5")
	provider.now = atHour(SnapdUpgradeHour)

	err := j.Run()

	assert.NoError(t, err)
	assert.Equal(t, 0, upgrader.called)
	assert.Empty(t, master.offered)
}

func TestSnapdUpgradeJob_SkipsOutsideWindow(t *testing.T) {
	j, upgrader, _, provider := newJob("2.60.0", "2.59.5")
	provider.now = atHour(23)

	err := j.Run()

	assert.NoError(t, err)
	assert.Equal(t, 0, upgrader.called)
}

func TestSnapdUpgradeJob_RunsOncePerDay(t *testing.T) {
	j, upgrader, _, provider := newJob("2.60.0", "2.59.5")

	provider.now = atHour(SnapdUpgradeHour)
	_ = j.Run()
	provider.now = atHour(SnapdUpgradeHour).Add(10 * time.Minute)
	_ = j.Run()

	assert.Equal(t, 1, upgrader.called)

	provider.now = atHour(SnapdUpgradeHour).AddDate(0, 0, 1)
	_ = j.Run()

	assert.Equal(t, 2, upgrader.called)
}

func TestSnapdUpgradeJob_RetriesWithinHourWhenBusy(t *testing.T) {
	j, upgrader, master, provider := newJob("2.60.0", "2.59.5")
	provider.now = atHour(SnapdUpgradeHour)
	master.busy = true

	err := j.Run()
	assert.Error(t, err)
	assert.Equal(t, 0, upgrader.called)

	master.busy = false
	provider.now = atHour(SnapdUpgradeHour).Add(5 * time.Minute)
	err = j.Run()

	assert.NoError(t, err)
	assert.Equal(t, 1, upgrader.called)
}
