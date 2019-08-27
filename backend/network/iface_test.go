package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIPv4(t *testing.T) {

	ip, err := LocalIp(false)

 assert.Nil(t, err)
 assert.NotNil(t, ip)
}

func TestIPv6(t *testing.T) {

	ip, err := LocalIp(true)

 assert.Nil(t, err)
 assert.NotNil(t, ip)
}
