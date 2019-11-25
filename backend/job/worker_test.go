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

type installerMock struct {
	installed int
}

func (mock *installerMock) Upgrade() {
	mock.installed++
}

type storageMock struct {
	formatted int
}

func (mock *storageMock) Format(device string) {
	mock.formatted++
}

func TestBackupCreate(t *testing.T) {
	master := &masterMock{}
	backup := &backupMock{}
	storage := &storageMock{}
	installer := &installerMock{}
	worker := NewWorker(master, backup, installer, storage)

	master.Offer(JobBackupCreate{"app"})
	worker.Do()

	assert.Equal(t, 1, backup.created)
	assert.Equal(t, 1, master.completed)

}

func TestNotSupported(t *testing.T) {
	master := &masterMock{}
	backup := &backupMock{}
	storage := &storageMock{}
	installer := &installerMock{}
	worker := NewWorker(master, backup, installer, storage)

	master.Offer("not supported type")
	worker.Do()

	assert.Equal(t, 1, master.completed)

}
