package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type backupMock struct {
	created int
}

func (mock *backupMock) List() ([]string, error) {
	return []string{"test"}, nil
}
func (mock *backupMock) Create(app string, file string) {
	mock.created++
}

type masterMock struct {
	job       interface{}
	taken     int
	completed int
}

func (mock *masterMock) Status() JobStatus {
	return JobStatusIdle
}
func (mock *masterMock) Offer(job interface{}) error {
	mock.job = job
	return nil
}
func (mock *masterMock) Take() (interface{}, error) {
	mock.taken++
	return mock.job, nil
}

func (mock *masterMock) Complete() error {
	mock.completed++
	return nil
}

func TestBackupCreate(t *testing.T) {
	master := &masterMock{}
	backup := &backupMock{}
	worker := NewWorker(master, backup)

	master.Offer(JobBackupCreate{"", ""})
	worker.Do()

	assert.Equal(t, 1, backup.created)
	assert.Equal(t, 1, master.completed)

}
