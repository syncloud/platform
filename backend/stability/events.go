package stability

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

const (
	maxLogFileBytes = 256 * 1024
	keepEvents      = 1000
	defaultLimit    = 100
)

type EventKind string

const (
	EventKindZramEnabled    EventKind = "zram_enabled"
	EventKindZramModuleLoad EventKind = "zram_module_loaded"
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
	App        string    `json:"app,omitempty"`
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
	if err := json.NewEncoder(f).Encode(e); err != nil {
		f.Close()
		return err
	}
	size := int64(0)
	if st, err := f.Stat(); err == nil {
		size = st.Size()
	}
	f.Close()
	if size > maxLogFileBytes {
		return l.rotateLocked(keepEvents)
	}
	return nil
}

func (l *EventLog) Recent(limit int) ([]Event, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.readLastLocked(limit, true)
}

func (l *EventLog) readLastLocked(limit int, reverse bool) ([]Event, error) {
	f, err := os.Open(l.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Event{}, nil
		}
		return nil, err
	}
	defer f.Close()
	ring := make([]Event, limit)
	n := 0
	dec := json.NewDecoder(f)
	for {
		var e Event
		if err := dec.Decode(&e); err != nil {
			break
		}
		ring[n%limit] = e
		n++
	}
	count := n
	if count > limit {
		count = limit
	}
	out := make([]Event, count)
	for i := 0; i < count; i++ {
		var idx int
		if reverse {
			idx = ((n - 1 - i) % limit + limit) % limit
		} else {
			start := 0
			if n > limit {
				start = n % limit
			}
			idx = (start + i) % limit
		}
		out[i] = ring[idx]
	}
	return out, nil
}

func (l *EventLog) rotateLocked(keep int) error {
	events, err := l.readLastLocked(keep, false)
	if err != nil {
		return err
	}
	tmp := l.path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	for _, e := range events {
		if err := enc.Encode(e); err != nil {
			f.Close()
			os.Remove(tmp)
			return err
		}
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, l.path)
}
