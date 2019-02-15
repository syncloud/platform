package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusIdle(t *testing.T) {
	master := NewMaster()

	assert.Equal(t, master.status, JobStatusIdle)
}

func TestStatusBusy(t *testing.T) {
	master := NewMaster()
	err := master.BackupCreateJob("nextcloud", "n.bkp")
	assert.Equal(t, err, nil)
	assert.Equal(t, master.Status(), JobStatusWaiting)
	master.Take()
	assert.Equal(t, master.Status(), JobStatusBusy)
	master.Complete()
	assert.Equal(t, master.Status(), JobStatusIdle)

}
