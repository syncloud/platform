package stability

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type EventKind string

const (
	EventKindZramEnabled    EventKind = "zram_enabled"
	EventKindSwapoffFile    EventKind = "swapoff_file"
	EventKindPressure       EventKind = "pressure_detected"
	EventKindVictimSigterm  EventKind = "victim_sigterm"
	EventKindVictimSigkill  EventKind = "victim_sigkill"
)

type Event struct {
	Time       time.Time `json:"time"`
	Kind       EventKind `json:"kind"`
	Message    string    `json:"message,omitempty"`
	PID        int       `json:"pid,omitempty"`
	Comm       string    `json:"comm,omitempty"`
	RSSkb      uint64    `json:"rss_kb,omitempty"`
	Cgroup     string    `json:"cgroup,omitempty"`
	AvailRatio float64   `json:"avail_ratio,omitempty"`
	PSIavg10   float64   `json:"psi_avg10,omitempty"`
	Path       string    `json:"path,omitempty"`
	SizeBytes  uint64    `json:"size_bytes,omitempty"`
}

type EventLog struct {
	path string
	mu   sync.Mutex
}

func NewEventLog(path string) *EventLog {
	return &EventLog{path: path}
}

func (l *EventLog) Append(e Event) error {
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(e)
}

func (l *EventLog) Recent(limit int) ([]Event, error) {
	if limit <= 0 {
		limit = 100
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.Open(l.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Event{}, nil
		}
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var all []Event
	for {
		var e Event
		if err := dec.Decode(&e); err != nil {
			break
		}
		all = append(all, e)
	}
	if len(all) > limit {
		all = all[len(all)-limit:]
	}
	for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
		all[i], all[j] = all[j], all[i]
	}
	return all, nil
}
