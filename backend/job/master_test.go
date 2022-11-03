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
	_, err := master.Take()
	assert.NotNil(t, err)
	assert.Equal(t, Idle, master.status)
}

func TestTakeWaiting(t *testing.T) {
	master := NewMaster()
	err := master.Offer("test", func() error { return nil })
	_, err = master.Take()
	assert.Nil(t, err)
	//assert.Equal(t, "job", job())
	assert.Equal(t, Busy, master.status)
}

func TestTakeBusy(t *testing.T) {
	master := NewMaster()
	_ = master.Offer("test", func() error { return nil })
	_, _ = master.Take()
	_, err := master.Take()
	assert.NotNil(t, err)
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

	_, err = master.Take()
	assert.Nil(t, err)
	assert.Equal(t, "test", master.Status().Name)

	err = master.Complete()
	assert.Nil(t, err)
	assert.Equal(t, "", master.Status().Name)
	assert.Equal(t, Idle, master.status)
}
