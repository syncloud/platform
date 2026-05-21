package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/stability"
)

func main() {
	logger := log.Default()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mem := stability.NewMemInfo("/proc")
	dataDir := os.Getenv("SNAP_DATA")
	if dataDir == "" {
		dataDir = "/var/snap/platform/current"
	}
	commonDir := os.Getenv("SNAP_COMMON")
	if commonDir == "" {
		commonDir = "/var/snap/platform/common"
	}
	eventsPath := dataDir + "/stability-events.jsonl"
	stability.MigrateEventLog(commonDir+"/stability-events.jsonl", eventsPath, logger)
	events := stability.NewEventLog(eventsPath)
	z := stability.NewZram(mem, stability.SwaponSyscall, stability.SwapoffSyscall, events, logger)
	if err := z.EnsureConfigured(); err != nil {
		logger.Sugar().Warnf("stability: zram setup failed (continuing): %v", err)
	}

	scan := stability.NewProcScanner("/proc")
	w := stability.NewWatcher(mem, scan, func(pid int, sig syscall.Signal) error {
		return syscall.Kill(pid, sig)
	}, events, logger)

	if err := w.Run(ctx); err != nil && err != context.Canceled {
		logger.Sugar().Errorf("stability: watcher exited: %v", err)
		os.Exit(1)
	}
}
