package event

import (
	"sync"

	"go.uber.org/zap"
)

type AccessChangeTrigger interface {
	RunAccessChangeEvent() error
}

type AccessChangeRunner struct {
	trigger AccessChangeTrigger
	logger  *zap.Logger
	spawn   func(func())
	mutex   sync.Mutex
	dirty   bool
	running bool
}

func NewAccessChangeRunner(trigger AccessChangeTrigger, logger *zap.Logger) *AccessChangeRunner {
	return &AccessChangeRunner{
		trigger: trigger,
		logger:  logger,
		spawn:   func(f func()) { go f() },
	}
}

func (r *AccessChangeRunner) Trigger() {
	r.mutex.Lock()
	r.dirty = true
	if r.running {
		r.mutex.Unlock()
		return
	}
	r.running = true
	r.mutex.Unlock()
	r.spawn(r.run)
}

func (r *AccessChangeRunner) run() {
	for {
		r.mutex.Lock()
		if !r.dirty {
			r.running = false
			r.mutex.Unlock()
			return
		}
		r.dirty = false
		r.mutex.Unlock()
		if err := r.trigger.RunAccessChangeEvent(); err != nil {
			r.logger.Error("access change event failed", zap.Error(err))
		}
	}
}
