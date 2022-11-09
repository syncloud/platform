package event

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/snap/model"
	"testing"
)

type SnapServerStub struct {
	snaps []model.Snap
}

func (s SnapServerStub) Snaps() ([]model.Snap, error) {
	return s.snaps, nil
}

type SnapCliStub struct {
	runs []string
}

func (e *SnapCliStub) Run(name string) error {
	e.runs = append(e.runs, name)
	return nil
}

func TestEvent_All(t *testing.T) {
	snapCli := &SnapCliStub{}
	snapd := &SnapServerStub{
		snaps: []model.Snap{
			{
				Name: "app1", Summary: "",
				Apps: []model.App{
					{Name: "event1", Snap: "app1"},
				},
			},
			{
				Name: "app2", Summary: "",
				Apps: []model.App{
					{Name: "event1", Snap: "app2"},
					{Name: "event2", Snap: "app2"},
				},
			},
		},
	}
	trigger := New(snapd, snapCli, log.Default())
	err := trigger.RunEventOnAllApps("event1")
	assert.Nil(t, err)
	assert.Len(t, snapCli.runs, 2)
	assert.Contains(t, snapCli.runs, "app1.event1")
	assert.Contains(t, snapCli.runs, "app2.event1")
}

func TestEvent_Filter(t *testing.T) {
	snapCli := &SnapCliStub{}
	snapd := &SnapServerStub{
		snaps: []model.Snap{
			{
				Name: "app1", Summary: "",
				Apps: []model.App{
					{Name: "event1", Snap: "app1"},
				},
			},
			{
				Name: "app2", Summary: "",
				Apps: []model.App{
					{Name: "event1", Snap: "app2"},
					{Name: "event2", Snap: "app2"},
				},
			},
		},
	}
	trigger := New(snapd, snapCli, log.Default())
	err := trigger.RunEventOnAllApps("event2")
	assert.Nil(t, err)
	assert.Len(t, snapCli.runs, 1)
	assert.Contains(t, snapCli.runs, "app2.event2")
}
