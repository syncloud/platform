package systemd

import (
	"github.com/syncloud/platform/cli"
	"strings"
)

type Journal struct {
	executor cli.Executor
}

func NewJournal(executor cli.Executor) *Journal {
	return &Journal{
		executor: executor,
	}
}

func (c *Journal) read(predicate func(string) bool, args ...string) []string {

	args = append(args, "-n", "1000", "--no-pager")
	output, err := c.executor.CombinedOutput("journalctl", args...)
	if err != nil {
		return []string{err.Error()}
	}
	var logs []string
	rawLogs := strings.Split(string(output), "\n")
	for _, line := range rawLogs {
		if predicate(line) {
			logs = append(logs, line)
		}
	}
	last := len(logs) - 1
	for i := 0; i < len(logs)/2; i++ {
		logs[i], logs[last-i] = logs[last-i], logs[i]
	}
	return logs
}

func (c *Journal) ReadAll(predicate func(string) bool) []string {
	return c.read(predicate)
}

func (c *Journal) ReadBackend(predicate func(string) bool) []string {
	return c.read(predicate, "-u", "snap.platform.backend")
}
