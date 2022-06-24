package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MasterStub struct {
	job       func()
	taken     int
	completed int
}

func (m *MasterStub) Take() (func(), error) {
	m.taken++
	return m.job, nil
}

func (m *MasterStub) Complete() error {
	m.completed++
	return nil
}

func TestJob(t *testing.T) {
	master := &MasterStub{}
	worker := NewWorker(master)

	ran := false
	master.job = func() { ran = true }
	worker.Do()

	assert.True(t, ran)
	assert.Equal(t, 1, master.completed)

}
