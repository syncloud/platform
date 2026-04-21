package timezone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/syncloud/platform/cli"
)

type Store interface {
	SetTimezone(tz string)
	GetTimezone() string
}

type Applier struct {
	executor cli.Executor
	store    Store
}

func NewApplier(executor cli.Executor, store Store) *Applier {
	return &Applier{executor: executor, store: store}
}

func (a *Applier) Apply(tz string) error {
	if !IsValid(tz) {
		return fmt.Errorf("invalid timezone: %s", tz)
	}
	output, err := a.executor.CombinedOutput("timedatectl", "set-timezone", tz)
	if err != nil {
		return fmt.Errorf("timedatectl set-timezone failed: %s: %w", string(output), err)
	}
	a.store.SetTimezone(tz)
	return nil
}

func (a *Applier) Current() string {
	return a.store.GetTimezone()
}

func IsValid(tz string) bool {
	if tz == "" || strings.Contains(tz, "..") || strings.HasPrefix(tz, "/") {
		return false
	}
	info, err := os.Stat(filepath.Join("/usr/share/zoneinfo", tz))
	if err != nil {
		return false
	}
	return !info.IsDir()
}
