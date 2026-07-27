package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type fakeAccessTrigger struct {
	runs   int
	onRun  func()
	failOn int
}

func (f *fakeAccessTrigger) RunAccessChangeEvent() error {
	f.runs++
	if f.onRun != nil {
		f.onRun()
	}
	return nil
}

func syncRunner(trigger AccessChangeTrigger) *AccessChangeRunner {
	r := NewAccessChangeRunner(trigger, zap.NewNop())
	r.spawn = func(f func()) { f() }
	return r
}

func TestAccessChangeRunner_RunsOnce(t *testing.T) {
	trigger := &fakeAccessTrigger{}
	r := syncRunner(trigger)
	r.Trigger()
	assert.Equal(t, 1, trigger.runs)
}

func TestAccessChangeRunner_ReRunsOnceWhenDirtiedMidRun(t *testing.T) {
	trigger := &fakeAccessTrigger{}
	var r *AccessChangeRunner
	trigger.onRun = func() {
		if trigger.runs == 1 {
			r.Trigger()
		}
	}
	r = syncRunner(trigger)
	r.Trigger()
	assert.Equal(t, 2, trigger.runs)
}

func TestAccessChangeRunner_CoalescesManyMidRunTriggers(t *testing.T) {
	trigger := &fakeAccessTrigger{}
	var r *AccessChangeRunner
	trigger.onRun = func() {
		if trigger.runs == 1 {
			r.Trigger()
			r.Trigger()
			r.Trigger()
		}
	}
	r = syncRunner(trigger)
	r.Trigger()
	assert.Equal(t, 2, trigger.runs)
}
