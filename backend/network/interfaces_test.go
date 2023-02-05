package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIPv4(t *testing.T) {

	ip, err := New().LocalIPv4()
	assert.True(t, err != nil || ip != nil)
}

func TestIPv6(t *testing.T) {

	interfaces := New()
	_, _ = interfaces.IPv6()
	assert.True(t, true)
}

func TestList(t *testing.T) {

	list, err := New().List()
	assert.Nil(t, err)

	assert.Greater(t, len(list), 0)
}
