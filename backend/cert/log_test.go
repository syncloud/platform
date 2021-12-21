package cert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type JournalCtlStub struct {
	logs string
}

func (j *JournalCtlStub) Read(unit string) (string, error) {
	return j.logs, nil
}

func TestLog(t *testing.T) {

	journalCtl := &JournalCtlStub{
		logs: `
3 log {"category": "certificate"}
2 log
1 log {"category": "certificate", "key": "value"}
`,
	}
	reader := NewReader(journalCtl)
	logs := reader.Read()
	assert.Equal(t, []string{
		`1 log {"category": "certificate", "key": "value"}`,
		`3 log {"category": "certificate"}`,
	}, logs)
}
