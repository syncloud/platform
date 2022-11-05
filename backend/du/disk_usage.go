package du

import (
	"github.com/syncloud/platform/cli"
	"regexp"
	"strconv"
)

type DiskUsage interface {
	Used(path string) (uint64, error)
}

type ShellDiskUsage struct {
	executor cli.CommandExecutor
}

func New(executor cli.CommandExecutor) *ShellDiskUsage {
	return &ShellDiskUsage{executor}
}

func (d *ShellDiskUsage) Used(path string) (uint64, error) {
	out, err := d.executor.CommandOutput("du", "-s", path)
	if err != nil {
		return 0, err
	}
	r := *regexp.MustCompile(`(\d+).*`)
	match := r.FindStringSubmatch(string(out))
	i, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}
