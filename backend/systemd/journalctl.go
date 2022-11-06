package systemd

import (
	"github.com/syncloud/platform/cli"
	"strings"
)

type JournalCtl struct {
	executor cli.Executor
}

type JournalCtlReader interface {
	ReadAll(predicate func(string) bool) []string
	ReadBackend(predicate func(string) bool) []string
}

func NewJournalCtl(executor cli.Executor) *JournalCtl {
	return &JournalCtl{
		executor: executor,
	}
}

func (c *JournalCtl) read(predicate func(string) bool, args ...string) []string {

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

func (c *JournalCtl) ReadAll(predicate func(string) bool) []string {
	return c.read(predicate)
}

func (c *JournalCtl) ReadBackend(predicate func(string) bool) []string {
	return c.read(predicate, "-u", "snap.platform.backend")
}
