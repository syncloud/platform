package systemd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type JournalCtlExecutorStub struct {
	output string
}

func (e *JournalCtlExecutorStub) CombinedOutput(_ string, _ ...string) ([]byte, error) {
	return []byte(e.output), nil
}

func Test_ReadBackend(t *testing.T) {

	journalCtl := &JournalCtlExecutorStub{
		output: `
3 log {"category": "cat1"}
2 log
1 log {"category": "cat1", "key": "value"}
`,
	}
	reader := NewJournal(journalCtl)
	logs := reader.ReadBackend(func(line string) bool {
		return strings.Contains(line, fmt.Sprintf(`"%s": "%s"`, "category", "cat1"))
	})
	assert.Equal(t, []string{
		`1 log {"category": "cat1", "key": "value"}`,
		`3 log {"category": "cat1"}`,
	}, logs)
}
