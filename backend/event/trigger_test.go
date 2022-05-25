package event

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/snap/model"
	"strings"
	"testing"
)

type SnapdStub struct {
	snaps []model.Snap
}

func (s SnapdStub) InstalledSnaps() ([]model.Snap, error) {
	return s.snaps, nil
}

type ExecutorStub struct {
	executions []string
}

func (e *ExecutorStub) CommandOutput(name string, arg ...string) ([]byte, error) {
	e.executions = append(e.executions, fmt.Sprintf("%s %s", name, strings.Join(arg, " ")))
	return make([]byte, 0), nil
}

func TestEvent_All(t *testing.T) {
	executor := &ExecutorStub{}
	snapd := &SnapdStub{
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
	trigger := New(snapd, executor)
	err := trigger.RunEventOnAllApps("event1")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 2)
	assert.Contains(t, executor.executions, "snap run app1.event1")
	assert.Contains(t, executor.executions, "snap run app2.event1")
}

func TestEvent_Filter(t *testing.T) {
	executor := &ExecutorStub{}
	snapd := &SnapdStub{
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
	trigger := New(snapd, executor)
	err := trigger.RunEventOnAllApps("event2")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 1)
	assert.Contains(t, executor.executions, "snap run app2.event2")
}
