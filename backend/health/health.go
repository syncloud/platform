package health

import (
	"github.com/syncloud/platform/stability"
)

type Health struct {
	events    *stability.EventLog
	collector *Collector
}

func NewHealth(events *stability.EventLog, collector *Collector) *Health {
	return &Health{events: events, collector: collector}
}

func (h *Health) Events(limit int) ([]stability.Event, error) {
	return h.events.Recent(limit)
}

func (h *Health) Metrics() (Snapshot, error) {
	return h.collector.Snapshot()
}
