package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusIdle(t *testing.T) {
	master := NewMaster()
	assert.Equal(t, master.status, Idle)
}

func TestOfferIdle(t *testing.T) {
	master := NewMaster()
	err := master.Offer("test", func() error { return nil })
	assert.Nil(t, err)
	assert.Equal(t, Waiting, master.status)
}

func TestOfferBusy(t *testing.T) {
	master := NewMaster()
	_ = master.Offer("test", func() error { return nil })
	err := master.Offer("test", func() error { return nil })
	assert.NotNil(t, err)
	assert.Equal(t, Waiting, master.status)
}

func TestTakeIdle(t *testing.T) {
	master := NewMaster()
	job := master.Take()
	assert.Nil(t, job)
	assert.Equal(t, Idle, master.status)
}

func TestTakeWaiting(t *testing.T) {
	master := NewMaster()
	err := master.Offer("test", func() error { return nil })
	assert.Nil(t, err)
	job := master.Take()
	assert.NotNil(t, job)
	//assert.Equal(t, "job", job())
	assert.Equal(t, Busy, master.status)
}

func TestTakeBusy(t *testing.T) {
	master := NewMaster()
	err := master.Offer("test", func() error { return nil })
	assert.Nil(t, err)
	job := master.Take()
	assert.NotNil(t, job)
	job = master.Take()
	assert.Nil(t, job)
	assert.Equal(t, Busy, master.status)
}

func TestCompleteIdle(t *testing.T) {
	master := NewMaster()
	err := master.Complete()
	assert.NotNil(t, err)
	assert.Equal(t, Idle, master.status)
}

func TestCompleteWaiting(t *testing.T) {
	master := NewMaster()
	_ = master.Offer("test", func() error { return nil })
	err := master.Complete()
	assert.NotNil(t, err)
	assert.Equal(t, Waiting, master.status)
}

func TestCompleteBusy(t *testing.T) {
	master := NewMaster()
	err := master.Offer("test", func() error { return nil })
	assert.Nil(t, err)
	assert.Equal(t, "test", master.Status().Name)

	job := master.Take()
	assert.NotNil(t, job)
	assert.Equal(t, "test", master.Status().Name)

	err = master.Complete()
	assert.Nil(t, err)
	assert.Equal(t, "", master.Status().Name)
	assert.Equal(t, Idle, master.status)
}
