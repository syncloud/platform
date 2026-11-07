package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

func TestAlerts_NoProblems(t *testing.T) {
	passed := true
	health := &HealthGetterStub{health: model.StorageHealth{
		Disks: []model.DiskHealth{{Device: "/dev/sda", Smart: model.SmartHealth{Available: true, Passed: &passed}}},
		Scrub: &model.ScrubStatus{State: model.ScrubFinished},
	}}
	alerts, err := NewAlerts(health, &AvailableDisksStub{disks: []model.Disk{{Device: "/dev/sda"}}}, &AlertsConfigStub{}).Get()
	assert.Nil(t, err)
	assert.Equal(t, &model.StorageAlerts{Problems: []model.StorageProblem{}}, alerts)
}

func TestAlerts_NewUntilSeen(t *testing.T) {
	passed := true
	failed := false
	health := &HealthGetterStub{health: model.StorageHealth{
		Disks: []model.DiskHealth{
			{Device: "/dev/sda", Smart: model.SmartHealth{Available: true, Passed: &passed, SelfTest: model.SmartTestResult{LastPassed: &failed}}},
			{Device: "/dev/sdb", Smart: model.SmartHealth{Available: true, Passed: &failed}, Io: model.KernelEvents{Count: 3}},
		},
		Scrub: &model.ScrubStatus{State: model.ScrubFinished, Errors: "read=12 csum=3", Uncorrectable: 15},
	}}
	disks := &AvailableDisksStub{disks: []model.Disk{{Device: "/dev/sda"}, {Device: "/dev/sdb", HasErrors: true}}}
	config := &AlertsConfigStub{}
	alerts := NewAlerts(health, disks, config)

	result, err := alerts.Get()
	assert.Nil(t, err)
	assert.True(t, result.New)
	assert.Equal(t, []model.StorageProblem{
		{Kind: model.ProblemErrorCounters, Device: "/dev/sdb"},
		{Kind: model.ProblemIoErrors, Device: "/dev/sdb"},
		{Kind: model.ProblemScrubErrors},
		{Kind: model.ProblemSelfTestFailed, Device: "/dev/sda"},
		{Kind: model.ProblemSmartFailed, Device: "/dev/sdb"},
	}, result.Problems)

	assert.Nil(t, alerts.Seen())
	result, err = alerts.Get()
	assert.Nil(t, err)
	assert.False(t, result.New)
	assert.Len(t, result.Problems, 5)

	health.health.Disks[0].Io = model.KernelEvents{Count: 1}
	result, err = alerts.Get()
	assert.Nil(t, err)
	assert.True(t, result.New)
}
