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
	err := master.Offer("test", func() {})
	assert.Nil(t, err)
	assert.Equal(t, Waiting, master.status)
}

func TestOfferBusy(t *testing.T) {
	master := NewMaster()
	_ = master.Offer("test", func() {})
	err := master.Offer("test", func() {})
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
	err := master.Offer("test", func() {})
	_, err = master.Take()
	assert.Nil(t, err)
	//assert.Equal(t, "job", job())
	assert.Equal(t, Busy, master.status)
}

func TestTakeBusy(t *testing.T) {
	master := NewMaster()
	_ = master.Offer("test", func() {})
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
	_ = master.Offer("test", func() {})
	err := master.Complete()
	assert.NotNil(t, err)
	assert.Equal(t, Waiting, master.status)
}

func TestCompleteBusy(t *testing.T) {
	master := NewMaster()
	err := master.Offer("test", func() {})
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
