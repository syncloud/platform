package main

import (
	"path"
	"syscall"

	"github.com/syncloud/platform/hook"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/stability"
)

func main() {
	logger := log.Default()
	mem := stability.NewMemInfo("/proc")
	eventsPath := path.Join(hook.DataDir, "stability-events.jsonl")
	stability.MigrateEventLog(path.Join(hook.CommonDir, "stability-events.jsonl"), eventsPath, logger)
	events := stability.NewEventLog(eventsPath)
	zram := stability.NewZram(mem, stability.SwaponSyscall, stability.SwapoffSyscall, stability.ModprobeLoad, events, logger)
	if err := zram.EnsureConfigured(); err != nil {
		logger.Sugar().Warnf("stability: zram setup failed (continuing): %v", err)
	}

	scanner := stability.NewProcScanner("/proc")
	watcher := stability.NewWatcher(mem, scanner, func(pid int, sig syscall.Signal) error {
		return syscall.Kill(pid, sig)
	}, events, logger)

	watcher.Run()
}
