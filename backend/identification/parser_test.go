package identification

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMac(t *testing.T) {
	mac, err := GetMac()
	assert.Nil(t, err)
	assert.NotEmpty(t, mac)
}
