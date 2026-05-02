package nginx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomProxy_ReloadCustomProxy_TriggersSystemd(t *testing.T) {
	nginx, systemd, _ := newTestNginx(t, "example.com", []ProxyEntry{
		{Name: "myapp", Host: "192.168.1.10", Port: 8080},
	})
	assert.NoError(t, nginx.ReloadCustomProxy())
	assert.Equal(t, "platform.nginx-custom-proxy", systemd.reloadedService)
}
