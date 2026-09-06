package btrfs

import (
	"fmt"
	"strings"

	"github.com/syncloud/platform/cli"
)

type ModuleLoader struct {
	executor cli.Executor
}

func NewModuleLoader(executor cli.Executor) *ModuleLoader {
	return &ModuleLoader{executor: executor}
}

func (l *ModuleLoader) Load(name string) error {
	output, err := l.executor.CombinedOutput("modprobe", name)
	if err != nil {
		return fmt.Errorf("modprobe %s: %w: %s", name, err, strings.TrimSpace(string(output)))
	}
	return nil
}
