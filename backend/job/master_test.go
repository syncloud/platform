package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusIdle(t *testing.T) {
	master := NewMaster()
	assert.Equal(t, master.status, JobStatusIdle)
}

func TestOfferIdle(t *testing.T) {
	master := NewMaster()
	err := master.Offer(func() {})
	assert.Nil(t, err)
	assert.Equal(t, JobStatusWaiting, master.Status())
}

func TestOfferBusy(t *testing.T) {
	master := NewMaster()
	_ = master.Offer(func() {})
	err := master.Offer(func() {})
	assert.NotNil(t, err)
	assert.Equal(t, JobStatusWaiting, master.Status())
}

func TestTakeIdle(t *testing.T) {
	master := NewMaster()
	_, err := master.Take()
	assert.NotNil(t, err)
	assert.Equal(t, JobStatusIdle, master.Status())
}

func TestTakeWaiting(t *testing.T) {
	master := NewMaster()
	err := master.Offer(func() {})
	_, err = master.Take()
	assert.Nil(t, err)
	//assert.Equal(t, "job", job())
	assert.Equal(t, JobStatusBusy, master.Status())
}

func TestTakeBusy(t *testing.T) {
	master := NewMaster()
	_ = master.Offer(func() {})
	_, _ = master.Take()
	_, err := master.Take()
	assert.NotNil(t, err)
	assert.Equal(t, JobStatusBusy, master.Status())
}

func TestCompleteIdle(t *testing.T) {
	master := NewMaster()
	err := master.Complete()
	assert.NotNil(t, err)
	assert.Equal(t, JobStatusIdle, master.Status())
}

func TestCompleteWaiting(t *testing.T) {
	master := NewMaster()
	_ = master.Offer(func() {})
	err := master.Complete()
	assert.NotNil(t, err)
	assert.Equal(t, JobStatusWaiting, master.Status())
}

func TestCompleteBusy(t *testing.T) {
	master := NewMaster()
	_ = master.Offer(func() {})
	_, _ = master.Take()
	err := master.Complete()
	assert.Nil(t, err)
	assert.Equal(t, JobStatusIdle, master.Status())
}
