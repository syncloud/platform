package kernel

import (
	"github.com/syncloud/platform/cli"
	"regexp"
	"strings"
)

var timestamped = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})\S*\s+(.*)$`)

type Log struct {
	executor cli.Executor
}

func NewLog(executor cli.Executor) *Log {
	return &Log{executor: executor}
}

func (l *Log) Events() (*Events, error) {
	output, err := l.executor.CombinedOutput("dmesg", "--time-format", "iso")
	if err != nil {
		return nil, err
	}
	return NewEvents(ParseLines(string(output))), nil
}

func ParseLines(output string) []Line {
	var lines []Line
	for _, raw := range strings.Split(output, "\n") {
		match := timestamped.FindStringSubmatch(strings.TrimSpace(raw))
		if match == nil {
			continue
		}
		lines = append(lines, Line{Time: match[1], Text: match[2]})
	}
	return lines
}
