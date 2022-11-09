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
	executor cli.Executor
}

func New(executor cli.Executor) *ShellDiskUsage {
	return &ShellDiskUsage{executor}
}

func (d *ShellDiskUsage) Used(path string) (uint64, error) {
	out, err := d.executor.CombinedOutput("du", "-s", path)
	if err != nil {
		return 0, err
	}
	r, err := regexp.Compile(`(\d+).*`)
	if err != nil {
		return 0, err
	}
	match := r.FindStringSubmatch(string(out))
	i, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}
