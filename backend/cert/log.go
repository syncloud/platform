package cert

import (
	"fmt"
	"github.com/syncloud/platform/log"
	"strings"
)

type Reader struct {
	journalCtl JournalCtl
}

func NewReader(journalCtl JournalCtl) *Reader {
	return &Reader{journalCtl: journalCtl}
}

func (l *Reader) Read() []string {

	output, err := l.journalCtl.Read("platform.backend")
	if err != nil {
		return []string{err.Error()}
	}

	var logs []string
	rawLogs := strings.Split(output, "\n")
	for _, line := range rawLogs {
		if strings.Contains(line, fmt.Sprintf(`"%s": "%s"`, log.CategoryKey, log.CategoryValue)) {
			logs = append(logs, line)
		}
	}
	last := len(logs) - 1
	for i := 0; i < len(logs)/2; i++ {
		logs[i], logs[last-i] = logs[last-i], logs[i]
	}
	return logs
}
