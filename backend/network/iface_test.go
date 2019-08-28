package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIPv4(t *testing.T) {

	ip, err := LocalIp(false)

 assert.True(t, err != nil || ip != nil)
}

func TestIPv6(t *testing.T) {

	ip, err := LocalIp(true)

 assert.True(t, err != nil || ip != nil)
}
