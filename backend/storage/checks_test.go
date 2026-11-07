package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

func newChecks(disks []model.Disk, smart *SmartStub, scrub *ScrubStub, btrfs *BtrfsStub, guard *GuardStub) *Checks {
	return NewChecks(&LsblkDisksStub{disks: disks}, smart, scrub, btrfs, &ScheduleConfigStub{}, guard, log.Default())
}

func TestChecks_StartScrub(t *testing.T) {
	scrub := &ScrubStub{status: model.ScrubStatus{State: model.ScrubFinished}}
	err := newChecks(nil, &SmartStub{}, scrub, &BtrfsStub{mounted: true}, &GuardStub{}).StartScrub()
	assert.Nil(t, err)
	assert.True(t, scrub.started)
}

func TestChecks_StartScrub_NoArray(t *testing.T) {
	scrub := &ScrubStub{}
	err := newChecks(nil, &SmartStub{}, scrub, &BtrfsStub{}, &GuardStub{}).StartScrub()
	assert.EqualError(t, err, "no active btrfs disk array")
	assert.False(t, scrub.started)
}

func TestChecks_StartScrub_RefusedDuringRefresh(t *testing.T) {
	scrub := &ScrubStub{}
	err := newChecks(nil, &SmartStub{}, scrub, &BtrfsStub{mounted: true}, &GuardStub{err: fmt.Errorf("an app install or update is in progress")}).StartScrub()
	assert.EqualError(t, err, "an app install or update is in progress")
	assert.False(t, scrub.started)
}

func TestChecks_StartScrub_AlreadyRunning(t *testing.T) {
	scrub := &ScrubStub{status: model.ScrubStatus{State: model.ScrubRunning}}
	err := newChecks(nil, &SmartStub{}, scrub, &BtrfsStub{mounted: true}, &GuardStub{}).StartScrub()
	assert.EqualError(t, err, "a scrub is already running")
	assert.False(t, scrub.started)
}

func TestChecks_StartSmartTest(t *testing.T) {
	smart := &SmartStub{}
	checks := newChecks([]model.Disk{{Device: "/dev/sda"}}, smart, &ScrubStub{}, &BtrfsStub{}, &GuardStub{})
	assert.Nil(t, checks.StartSmartTest("/dev/sda"))
	assert.EqualError(t, checks.StartSmartTest("/dev/mmcblk0"), "unknown disk: /dev/mmcblk0")
	assert.Equal(t, []string{"/dev/sda"}, smart.started)
}

func TestChecks_StartSmartTest_RefusedDuringRefresh(t *testing.T) {
	smart := &SmartStub{}
	checks := newChecks([]model.Disk{{Device: "/dev/sda"}}, smart, &ScrubStub{}, &BtrfsStub{}, &GuardStub{err: fmt.Errorf("busy")})
	assert.EqualError(t, checks.StartSmartTest("/dev/sda"), "busy")
	assert.Empty(t, smart.started)
}

func TestChecks_StartSmartTests(t *testing.T) {
	disks := []model.Disk{{Device: "/dev/sda"}, {Device: "/dev/sdb"}, {Device: "/dev/sdc"}, {Device: "/dev/sdd"}}
	smart := &SmartStub{
		health: map[string]model.SmartHealth{
			"/dev/sda": {Available: true},
			"/dev/sdb": {Available: false},
			"/dev/sdc": {Available: true, SelfTest: model.SmartTestResult{Running: true}},
			"/dev/sdd": {Available: true},
		},
		fail: map[string]bool{"/dev/sdd": true},
	}
	err := newChecks(disks, smart, &ScrubStub{}, &BtrfsStub{}, &GuardStub{}).StartSmartTests()
	assert.EqualError(t, err, "/dev/sdd: test refused")
	assert.Equal(t, []string{"/dev/sda"}, smart.started)
}

func TestChecks_ResetErrors(t *testing.T) {
	btrfs := &BtrfsStub{mounted: true}
	assert.Nil(t, newChecks(nil, &SmartStub{}, &ScrubStub{}, btrfs, &GuardStub{}).ResetErrors())
	assert.True(t, btrfs.reset)

	unmounted := &BtrfsStub{}
	assert.EqualError(t, newChecks(nil, &SmartStub{}, &ScrubStub{}, unmounted, &GuardStub{}).ResetErrors(), "no active btrfs disk array")
	assert.False(t, unmounted.reset)
}
