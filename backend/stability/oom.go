package stability

import (
	"context"
	"errors"
	"os"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type KillFn func(pid int, sig syscall.Signal) error

type Watcher struct {
	mem      *MemInfo
	scan     *ProcScanner
	protect  Protect
	kill     KillFn
	events   *EventLog
	log      *zap.Logger
	interval time.Duration
	availMin float64
	psiMax   float64
	grace    time.Duration
	selfPID  int
}

func NewWatcher(mem *MemInfo, scan *ProcScanner, kill KillFn, events *EventLog, log *zap.Logger) *Watcher {
	return &Watcher{
		mem:      mem,
		scan:     scan,
		protect:  DefaultProtect(),
		kill:     kill,
		events:   events,
		log:      log,
		interval: 2 * time.Second,
		availMin: 0.08,
		psiMax:   40.0,
		grace:    4 * time.Second,
		selfPID:  os.Getpid(),
	}
}

func (w *Watcher) Run(ctx context.Context) error {
	t := time.NewTicker(w.interval)
	defer t.Stop()
	w.log.Info("oom-watcher: started",
		zap.Duration("interval", w.interval),
		zap.Float64("avail_min", w.availMin),
		zap.Float64("psi_max", w.psiMax),
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			if err := w.tick(); err != nil {
				w.log.Warn("oom-watcher: tick error", zap.Error(err))
			}
		}
	}
}

func (w *Watcher) tick() error {
	snap, err := w.mem.Snapshot()
	if err != nil {
		return err
	}
	avail := snap.AvailableRatio()
	psi := 0.0
	psiOK := false
	if w.mem.PSIAvailable() {
		if v, err := w.mem.PSIMemoryAvg10(); err == nil {
			psi = v
			psiOK = true
		}
	}
	if !w.pressureExceeded(avail, psi, psiOK) {
		return nil
	}
	w.log.Warn("oom-watcher: pressure detected",
		zap.Float64("avail_ratio", avail),
		zap.Float64("psi_avg10", psi),
		zap.Bool("psi_ok", psiOK),
	)
	if w.events != nil {
		_ = w.events.Append(Event{Kind: EventKindPressure, AvailRatio: avail, PSIavg10: psi})
	}
	return w.killWorst()
}

func (w *Watcher) pressureExceeded(avail, psi float64, psiOK bool) bool {
	if avail < w.availMin {
		return true
	}
	if psiOK && psi > w.psiMax {
		return true
	}
	return false
}

func (w *Watcher) killWorst() error {
	cands, err := w.scan.Candidates(w.protect, w.selfPID)
	if err != nil {
		return err
	}
	if len(cands) == 0 {
		return ErrNoVictim
	}
	v := cands[0]
	w.log.Warn("oom-watcher: SIGTERM victim",
		zap.Int("pid", v.PID),
		zap.String("comm", v.Comm),
		zap.Uint64("rss_kb", v.RSSkB),
		zap.String("cgroup", v.Cgroup),
	)
	if w.events != nil {
		_ = w.events.Append(Event{Kind: EventKindVictimSigterm, PID: v.PID, Comm: v.Comm, RSSkb: v.RSSkB, Cgroup: v.Cgroup})
	}
	if err := w.kill(v.PID, syscall.SIGTERM); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return err
	}
	deadline := time.Now().Add(w.grace)
	for time.Now().Before(deadline) {
		if w.kill(v.PID, 0) != nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	w.log.Warn("oom-watcher: SIGKILL victim", zap.Int("pid", v.PID), zap.String("comm", v.Comm))
	if w.events != nil {
		_ = w.events.Append(Event{Kind: EventKindVictimSigkill, PID: v.PID, Comm: v.Comm, RSSkb: v.RSSkB, Cgroup: v.Cgroup})
	}
	if err := w.kill(v.PID, syscall.SIGKILL); err != nil && !errors.Is(err, syscall.ESRCH) {
		return err
	}
	return nil
}
