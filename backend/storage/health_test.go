package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

func TestHealth_Get(t *testing.T) {
	passed := true
	disks := []model.Disk{{Device: "/dev/sda"}, {Device: "/dev/sdb"}}
	smart := &SmartStub{health: map[string]model.SmartHealth{
		"/dev/sda": {Available: true, Passed: &passed},
		"/dev/sdb": {Message: "SMART is not available for this disk"},
	}}
	kernelLog := &KernelStub{dmesg: `2026-10-30T14:12:09,511402+00:00 ata2: hard resetting link
2026-10-31T02:01:45,000100+00:00 blk_update_request: I/O error, dev sdb, sector 123456 op 0x0:(READ)
`}
	ports := &PortsStub{ports: map[string]string{"/dev/sda": "ata1", "/dev/sdb": "ata2"}}
	scrub := &ScrubStub{status: model.ScrubStatus{State: model.ScrubFinished}}
	schedule := &ScheduleConfigStub{schedule: model.CheckSchedule{Scrub: true}}
	health := NewHealth(&LsblkDisksStub{disks: disks}, smart, kernelLog, ports, scrub, &BtrfsStub{mounted: true}, schedule, &GuardStub{}, log.Default())

	result, err := health.Get()
	assert.Nil(t, err)
	assert.Equal(t, []model.DiskHealth{
		{Device: "/dev/sda", Port: "ata1", Smart: model.SmartHealth{Available: true, Passed: &passed}},
		{Device: "/dev/sdb", Port: "ata2", Smart: model.SmartHealth{Message: "SMART is not available for this disk"},
			Link: model.KernelEvents{Count: 1, Last: "2026-10-30T14:12:09"},
			Io:   model.KernelEvents{Count: 1, Last: "2026-10-31T02:01:45"}},
	}, result.Disks)
	assert.Equal(t, &model.ScrubStatus{State: model.ScrubFinished}, result.Scrub)
	assert.Equal(t, model.CheckSchedule{Scrub: true}, result.Schedule)
	assert.False(t, result.Busy)
}

func TestHealth_Get_NoArrayNoScrub(t *testing.T) {
	health := NewHealth(&LsblkDisksStub{}, &SmartStub{}, &KernelStub{}, &PortsStub{}, &ScrubStub{}, &BtrfsStub{}, &ScheduleConfigStub{}, &GuardStub{err: fmt.Errorf("busy")}, log.Default())

	result, err := health.Get()
	assert.Nil(t, err)
	assert.Equal(t, []model.DiskHealth{}, result.Disks)
	assert.Nil(t, result.Scrub)
	assert.True(t, result.Busy)
}
