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
	commonDir := os.Getenv("SNAP_COMMON")
	if commonDir == "" {
		commonDir = "/var/snap/platform/common"
	}
	events := stability.NewEventLog(commonDir + "/stability-events.jsonl")
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
