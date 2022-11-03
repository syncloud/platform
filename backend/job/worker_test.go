package job

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
)

type MasterStub struct {
	job       func() error
	taken     int
	completed int
}

func (m *MasterStub) Take() (func() error, error) {
	m.taken++
	return m.job, nil
}

func (m *MasterStub) Complete() error {
	m.completed++
	return nil
}

func TestJob(t *testing.T) {
	master := &MasterStub{}
	worker := NewWorker(master, log.Default())

	ran := false
	master.job = func() error {
		ran = true
		return nil
	}
	worker.Do()

	assert.True(t, ran)
	assert.Equal(t, 1, master.completed)

}
