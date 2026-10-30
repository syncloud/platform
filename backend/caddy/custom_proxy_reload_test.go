package caddy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_ReloadCustomProxy_TriggersSystemd(t *testing.T) {
	caddy, systemd, _ := newTestCaddy(t, "example.com", []ProxyEntry{
		{Name: "myapp", Host: "192.168.1.10", Port: 8080},
	})
	assert.NoError(t, caddy.ReloadCustomProxy())
	assert.Equal(t, Service, systemd.reloadedService)
}
