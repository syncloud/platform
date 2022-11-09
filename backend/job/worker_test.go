package job

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"testing"
)

type MasterStub struct {
	job       Job
	taken     int
	completed int
}

func (m *MasterStub) Take() Job {
	m.taken++
	return m.job
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

func TestNoJob(t *testing.T) {
	master := &MasterStub{}
	worker := NewWorker(master, log.Default())

	master.job = nil
	worker.Do()

	assert.Equal(t, 0, master.completed)
}
