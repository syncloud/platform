package job

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/backup"
	"testing"
)

type backupMock struct {
	created  int
	restored int
}

func (mock *backupMock) List() ([]backup.File, error) {
	return []backup.File{backup.File{"dir", "test"}}, nil
}
func (mock *backupMock) Create(app string) {
	mock.created++
}
func (mock *backupMock) Restore(file string) {
	mock.restored++
}

type masterMock struct {
	job       func()
	taken     int
	completed int
}

func (mock *masterMock) Status() JobStatus {
	return JobStatusIdle
}
func (mock *masterMock) Offer(job func()) error {
	mock.job = job
	return nil
}
func (mock *masterMock) Take() (func(), error) {
	mock.taken++
	return mock.job, nil
}

func (mock *masterMock) Complete() error {
	mock.completed++
	return nil
}

type installerMock struct {
	installed int
}

func (mock *installerMock) Upgrade() {
	mock.installed++
}

type storageMock struct {
	formatted    int
	bootextended int
}

func (mock *storageMock) Format(device string) {
	mock.formatted++
}

func (mock *storageMock) BootExtend() {
	mock.bootextended++
}

func TestJob(t *testing.T) {
	master := &masterMock{}
	worker := NewWorker(master)

	ran := false
	err := master.Offer(func() { ran = true })
	worker.Do()

	assert.Nil(t, err)
	assert.True(t, ran)
	assert.Equal(t, 1, master.completed)

}
